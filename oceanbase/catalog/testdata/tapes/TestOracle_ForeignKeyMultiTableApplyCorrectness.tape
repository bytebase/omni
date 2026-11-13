-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_ForeignKeyMultiTableApplyCorrectness default#1
> exec
| DROP DATABASE IF EXISTS fkmulti_t_drop_column_and_its_fk
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_drop_column_and_its_fk DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_drop_column_and_its_fk
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_drop_column_and_its_fk' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_t_drop_column_and_its_fk.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti_t_drop_column_and_its_fk.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkmulti
< ok 2 0
> exec
| CREATE DATABASE fkmulti DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fkmulti`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `fk`
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP COLUMN `pid`
< ok 0 0
> query
| SHOW CREATE TABLE fkmulti.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkmulti_t_drop_referenced_table
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_drop_referenced_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_drop_referenced_table
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_drop_referenced_table' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| c
> query
| SHOW CREATE TABLE fkmulti_t_drop_referenced_table.c
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
| DROP DATABASE IF EXISTS fkmulti
< ok 2 0
> exec
| CREATE DATABASE fkmulti DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti.c
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
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fkmulti`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `fk`
< ok 0 0
> exec
| DROP TABLE `fkmulti`.`p`
< ok 0 0
> query
| SHOW CREATE TABLE fkmulti.c
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
| DROP DATABASE IF EXISTS fkmulti_t_drop_two_fks_sharing_index
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_drop_two_fks_sharing_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_drop_two_fks_sharing_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_drop_two_fks_sharing_index' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_t_drop_two_fks_sharing_index.c
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
> query
| SHOW CREATE TABLE fkmulti_t_drop_two_fks_sharing_index.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkmulti
< ok 1 0
> exec
| CREATE DATABASE fkmulti DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fka FOREIGN KEY (pid) REFERENCES p(id), CONSTRAINT fkb FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti.c
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
|   CONSTRAINT `fkb` FOREIGN KEY (`pid`) REFERENCES `fkmulti`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT ,
|   CONSTRAINT `fka` FOREIGN KEY (`pid`) REFERENCES `fkmulti`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `fka`
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `fkb`
< ok 0 0
> query
| SHOW CREATE TABLE fkmulti.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti.c
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
| DROP DATABASE IF EXISTS fkmulti_t_drop_one_of_two_fks_sharing_index
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_drop_one_of_two_fks_sharing_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_drop_one_of_two_fks_sharing_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fkb FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_drop_one_of_two_fks_sharing_index' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_t_drop_one_of_two_fks_sharing_index.c
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
|   CONSTRAINT `fkb` FOREIGN KEY (`pid`) REFERENCES `fkmulti_t_drop_one_of_two_fks_sharing_index`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti_t_drop_one_of_two_fks_sharing_index.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkmulti
< ok 2 0
> exec
| CREATE DATABASE fkmulti DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fka FOREIGN KEY (pid) REFERENCES p(id), CONSTRAINT fkb FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti.c
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
|   CONSTRAINT `fkb` FOREIGN KEY (`pid`) REFERENCES `fkmulti`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT ,
|   CONSTRAINT `fka` FOREIGN KEY (`pid`) REFERENCES `fkmulti`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `fka`
< ok 0 0
> query
| SHOW CREATE TABLE fkmulti.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti.c
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
|   CONSTRAINT `fkb` FOREIGN KEY (`pid`) REFERENCES `fkmulti`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkmulti_t_drop_fk_reuse_name_for_other_index
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_drop_fk_reuse_name_for_other_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_drop_fk_reuse_name_for_other_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, other_col INT, KEY fk (other_col))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_drop_fk_reuse_name_for_other_index' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_t_drop_fk_reuse_name_for_other_index.c
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
|   `other_col` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk` (`other_col`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti_t_drop_fk_reuse_name_for_other_index.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkmulti
< ok 2 0
> exec
| CREATE DATABASE fkmulti DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, other_col INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti.c
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
|   `other_col` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fkmulti`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `fk`
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` ADD KEY `fk` (`other_col`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE fkmulti.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti.c
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
|   `other_col` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk` (`other_col`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkmulti_t_drop_parent_and_child_name_order
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_drop_parent_and_child_name_order DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_drop_parent_and_child_name_order
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_drop_parent_and_child_name_order' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
> exec
| DROP DATABASE IF EXISTS fkmulti
< ok 3 0
> exec
| CREATE DATABASE fkmulti DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti
< ok 0 0
> exec
| CREATE TABLE a_parent (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE z_child (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES a_parent(id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| a_parent
< row
: bytes
| z_child
> query
| SHOW CREATE TABLE fkmulti.a_parent
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a_parent
: bytes
| CREATE TABLE `a_parent` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti.z_child
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| z_child
: bytes
| CREATE TABLE `z_child` (
|   `id` int(11) NOT NULL,
|   `pid` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fkmulti`.`a_parent`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `fkmulti`.`z_child` DROP FOREIGN KEY `fk`
< ok 0 0
> exec
| DROP TABLE `fkmulti`.`a_parent`
< ok 0 0
> exec
| DROP TABLE `fkmulti`.`z_child`
< ok 0 0
> exec
| DROP DATABASE IF EXISTS fkmulti_t_drop_two_unnamed_fks_sharing_index
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_drop_two_unnamed_fks_sharing_index DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_drop_two_unnamed_fks_sharing_index
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_drop_two_unnamed_fks_sharing_index' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_t_drop_two_unnamed_fks_sharing_index.c
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
> query
| SHOW CREATE TABLE fkmulti_t_drop_two_unnamed_fks_sharing_index.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkmulti
< ok 0 0
> exec
| CREATE DATABASE fkmulti DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES p(id), FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti.c
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
|   CONSTRAINT `c_OBFK_0000000000000001` FOREIGN KEY (`pid`) REFERENCES `fkmulti`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT ,
|   CONSTRAINT `c_OBFK_0000000000000002` FOREIGN KEY (`pid`) REFERENCES `fkmulti`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `c_OBFK_0000000000000001`
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `c_OBFK_0000000000000002`
< ok 0 0
> query
| SHOW CREATE TABLE fkmulti.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti.c
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
| DROP DATABASE IF EXISTS fkmulti_t_replace_primary_key_of_referenced_table
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_replace_primary_key_of_referenced_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_replace_primary_key_of_referenced_table
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL, u INT NOT NULL, PRIMARY KEY (id, u), UNIQUE KEY uu (u))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pu INT, CONSTRAINT fk FOREIGN KEY (pu) REFERENCES p(u))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_replace_primary_key_of_referenced_table' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_t_replace_primary_key_of_referenced_table.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL,
|   `pu` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk` FOREIGN KEY (`pu`) REFERENCES `fkmulti_t_replace_primary_key_of_referenced_table`.`p`(`u`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti_t_replace_primary_key_of_referenced_table.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   `u` int(11) NOT NULL,
|   PRIMARY KEY (`id`, `u`),
|   UNIQUE KEY `uu` (`u`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkmulti
< ok 2 0
> exec
| CREATE DATABASE fkmulti DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL, u INT NOT NULL, PRIMARY KEY (id), UNIQUE KEY uu (u))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pu INT, CONSTRAINT fk FOREIGN KEY (pu) REFERENCES p(u))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL,
|   `pu` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk` FOREIGN KEY (`pu`) REFERENCES `fkmulti`.`p`(`u`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   `u` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uu` (`u`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `fkmulti`.`p` DROP PRIMARY KEY, ADD PRIMARY KEY (`id`,`u`)
< ok 0 0
> query
| SHOW CREATE TABLE fkmulti.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   `u` int(11) NOT NULL,
|   PRIMARY KEY (`id`, `u`),
|   UNIQUE KEY `uu` (`u`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL,
|   `pu` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk` FOREIGN KEY (`pu`) REFERENCES `fkmulti`.`p`(`u`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
