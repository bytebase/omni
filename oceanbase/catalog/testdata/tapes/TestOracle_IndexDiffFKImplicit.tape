-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_IndexDiffFKImplicit default#1
> exec
| DROP DATABASE IF EXISTS idxfk_fk_named_implicit
< ok 0 0
> exec
| CREATE DATABASE idxfk_fk_named_implicit DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE idxfk_fk_named_implicit
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE idxfk_fk_named_implicit.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   `name` varchar(50) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE idxfk_fk_named_implicit.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `idxfk_fk_named_implicit`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS idxfk_fk_unnamed_implicit
< ok 0 0
> exec
| CREATE DATABASE idxfk_fk_unnamed_implicit DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE idxfk_fk_unnamed_implicit
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE idxfk_fk_unnamed_implicit.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   `name` varchar(50) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE idxfk_fk_unnamed_implicit.c
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
|   CONSTRAINT `c_OBFK_0000000000000001` FOREIGN KEY (`pid`) REFERENCES `idxfk_fk_unnamed_implicit`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS idxfk_fk_explicit_covering
< ok 0 0
> exec
| CREATE DATABASE idxfk_fk_explicit_covering DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE idxfk_fk_explicit_covering
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE idxfk_fk_explicit_covering.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   `name` varchar(50) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE idxfk_fk_explicit_covering.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `idxfk_fk_explicit_covering`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT ,
|   KEY `my_idx` (`pid`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS idxfk_fk_composite
< ok 0 0
> exec
| CREATE DATABASE idxfk_fk_composite DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE idxfk_fk_composite
< ok 0 0
> exec
| CREATE TABLE p (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a,b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT, CONSTRAINT fk FOREIGN KEY (x,y) REFERENCES p(a,b))
< ok 0 0
> query
| SHOW CREATE TABLE idxfk_fk_composite.p
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
|   PRIMARY KEY (`a`, `b`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE idxfk_fk_composite.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`x`, `y`) REFERENCES `idxfk_fk_composite`.`p`(`a`, `b`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS idxfk_fk_backing_name_collision
< ok 0 0
> exec
| CREATE DATABASE idxfk_fk_backing_name_collision DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE idxfk_fk_backing_name_collision
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, a INT, pid INT, KEY pid (a), FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE idxfk_fk_backing_name_collision.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   `name` varchar(50) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE idxfk_fk_backing_name_collision.c
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
|   CONSTRAINT `c_OBFK_0000000000000002` FOREIGN KEY (`pid`) REFERENCES `idxfk_fk_backing_name_collision`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT ,
|   KEY `pid` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS idxfk_fk_user_named_ibfk
< ok 0 0
> exec
| CREATE DATABASE idxfk_fk_user_named_ibfk DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE idxfk_fk_user_named_ibfk
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT c_ibfk_1 FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE idxfk_fk_user_named_ibfk.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   `name` varchar(50) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE idxfk_fk_user_named_ibfk.c
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
|   CONSTRAINT `c_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `idxfk_fk_user_named_ibfk`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
