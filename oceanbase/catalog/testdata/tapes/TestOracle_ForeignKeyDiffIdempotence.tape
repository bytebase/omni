-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_ForeignKeyDiffIdempotence default#1
> exec
| DROP DATABASE IF EXISTS fkdiff_s_named
< ok 0 0
> exec
| CREATE DATABASE fkdiff_s_named DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkdiff_s_named
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fkdiff_s_named.p
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
| SHOW CREATE TABLE fkdiff_s_named.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fkdiff_s_named`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkdiff_s_unnamed_autoname
< ok 0 0
> exec
| CREATE DATABASE fkdiff_s_unnamed_autoname DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkdiff_s_unnamed_autoname
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fkdiff_s_unnamed_autoname.p
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
| SHOW CREATE TABLE fkdiff_s_unnamed_autoname.c
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
|   CONSTRAINT `c_OBFK_0000000000000001` FOREIGN KEY (`pid`) REFERENCES `fkdiff_s_unnamed_autoname`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkdiff_s_user_named_ibfk
< ok 0 0
> exec
| CREATE DATABASE fkdiff_s_user_named_ibfk DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkdiff_s_user_named_ibfk
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT c_ibfk_1 FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fkdiff_s_user_named_ibfk.p
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
| SHOW CREATE TABLE fkdiff_s_user_named_ibfk.c
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
|   CONSTRAINT `c_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `fkdiff_s_user_named_ibfk`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkdiff_s_action_bare
< ok 0 0
> exec
| CREATE DATABASE fkdiff_s_action_bare DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkdiff_s_action_bare
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fkdiff_s_action_bare.p
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
| SHOW CREATE TABLE fkdiff_s_action_bare.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fkdiff_s_action_bare`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkdiff_s_action_restrict
< ok 0 0
> exec
| CREATE DATABASE fkdiff_s_action_restrict DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkdiff_s_action_restrict
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE RESTRICT ON UPDATE RESTRICT)
< ok 0 0
> query
| SHOW CREATE TABLE fkdiff_s_action_restrict.p
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
| SHOW CREATE TABLE fkdiff_s_action_restrict.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fkdiff_s_action_restrict`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkdiff_s_action_no_action
< ok 0 0
> exec
| CREATE DATABASE fkdiff_s_action_no_action DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkdiff_s_action_no_action
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE NO ACTION ON UPDATE NO ACTION)
< ok 0 0
> query
| SHOW CREATE TABLE fkdiff_s_action_no_action.p
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
| SHOW CREATE TABLE fkdiff_s_action_no_action.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fkdiff_s_action_no_action`.`p`(`id`) ON UPDATE NO ACTION ON DELETE NO ACTION 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkdiff_s_action_cascade_setnull
< ok 0 0
> exec
| CREATE DATABASE fkdiff_s_action_cascade_setnull DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkdiff_s_action_cascade_setnull
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE CASCADE ON UPDATE SET NULL)
< ok 0 0
> query
| SHOW CREATE TABLE fkdiff_s_action_cascade_setnull.p
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
| SHOW CREATE TABLE fkdiff_s_action_cascade_setnull.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fkdiff_s_action_cascade_setnull`.`p`(`id`) ON UPDATE SET NULL ON DELETE CASCADE 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkdiff_s_action_delete_only
< ok 0 0
> exec
| CREATE DATABASE fkdiff_s_action_delete_only DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkdiff_s_action_delete_only
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE CASCADE)
< ok 0 0
> query
| SHOW CREATE TABLE fkdiff_s_action_delete_only.p
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
| SHOW CREATE TABLE fkdiff_s_action_delete_only.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fkdiff_s_action_delete_only`.`p`(`id`) ON UPDATE RESTRICT ON DELETE CASCADE 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkdiff_s_action_update_only
< ok 0 0
> exec
| CREATE DATABASE fkdiff_s_action_update_only DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkdiff_s_action_update_only
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON UPDATE CASCADE)
< ok 0 0
> query
| SHOW CREATE TABLE fkdiff_s_action_update_only.p
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
| SHOW CREATE TABLE fkdiff_s_action_update_only.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fkdiff_s_action_update_only`.`p`(`id`) ON UPDATE CASCADE ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkdiff_s_action_setnull_both
< ok 0 0
> exec
| CREATE DATABASE fkdiff_s_action_setnull_both DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkdiff_s_action_setnull_both
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE SET NULL ON UPDATE SET NULL)
< ok 0 0
> query
| SHOW CREATE TABLE fkdiff_s_action_setnull_both.p
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
| SHOW CREATE TABLE fkdiff_s_action_setnull_both.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fkdiff_s_action_setnull_both`.`p`(`id`) ON UPDATE SET NULL ON DELETE SET NULL 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkdiff_s_composite
< ok 0 0
> exec
| CREATE DATABASE fkdiff_s_composite DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkdiff_s_composite
< ok 0 0
> exec
| CREATE TABLE p (a INT NOT NULL, b INT NOT NULL, c INT, PRIMARY KEY (a,b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT, CONSTRAINT fk FOREIGN KEY (x,y) REFERENCES p(a,b))
< ok 0 0
> query
| SHOW CREATE TABLE fkdiff_s_composite.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   PRIMARY KEY (`a`, `b`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkdiff_s_composite.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL,
|   `x` int(11) DEFAULT NULL,
|   `y` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk` FOREIGN KEY (`x`, `y`) REFERENCES `fkdiff_s_composite`.`p`(`a`, `b`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkdiff_s_explicit_covering_index
< ok 0 0
> exec
| CREATE DATABASE fkdiff_s_explicit_covering_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkdiff_s_explicit_covering_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fkdiff_s_explicit_covering_index.p
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
| SHOW CREATE TABLE fkdiff_s_explicit_covering_index.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fkdiff_s_explicit_covering_index`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT ,
|   KEY `my_idx` (`pid`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkdiff_s_backing_name_collision
< ok 0 0
> exec
| CREATE DATABASE fkdiff_s_backing_name_collision DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkdiff_s_backing_name_collision
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, a INT, pid INT, KEY pid (a), FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fkdiff_s_backing_name_collision.p
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
| SHOW CREATE TABLE fkdiff_s_backing_name_collision.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `pid` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `c_OBFK_0000000000000002` FOREIGN KEY (`pid`) REFERENCES `fkdiff_s_backing_name_collision`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT ,
|   KEY `pid` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkdiff_s_multiple_fks
< ok 0 0
> exec
| CREATE DATABASE fkdiff_s_multiple_fks DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkdiff_s_multiple_fks
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pa INT, pb INT, CONSTRAINT fka FOREIGN KEY (pa) REFERENCES p(a), CONSTRAINT fkb FOREIGN KEY (pb) REFERENCES p(b))
< ok 0 0
> query
| SHOW CREATE TABLE fkdiff_s_multiple_fks.p
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
| SHOW CREATE TABLE fkdiff_s_multiple_fks.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL,
|   `pa` int(11) DEFAULT NULL,
|   `pb` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fkb` FOREIGN KEY (`pb`) REFERENCES `fkdiff_s_multiple_fks`.`p`(`b`) ON UPDATE RESTRICT ON DELETE RESTRICT ,
|   CONSTRAINT `fka` FOREIGN KEY (`pa`) REFERENCES `fkdiff_s_multiple_fks`.`p`(`a`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
