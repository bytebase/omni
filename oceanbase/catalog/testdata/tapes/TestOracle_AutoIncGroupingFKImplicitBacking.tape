-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_AutoIncGroupingFKImplicitBacking default#1
> exec
| DROP DATABASE IF EXISTS fkmulti_t_aig_fk_implicit_backing
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_aig_fk_implicit_backing DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_aig_fk_implicit_backing
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT NOT NULL PRIMARY KEY, x INT, seq INT NOT NULL AUTO_INCREMENT, CONSTRAINT fk_seq FOREIGN KEY (seq) REFERENCES p (id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_aig_fk_implicit_backing' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_t_aig_fk_implicit_backing.c
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
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk_seq` FOREIGN KEY (`seq`) REFERENCES `fkmulti_t_aig_fk_implicit_backing`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti_t_aig_fk_implicit_backing.p
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
| CREATE TABLE c (id INT NOT NULL PRIMARY KEY, x INT)
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
|   `x` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
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
| ALTER TABLE `fkmulti`.`c` ADD COLUMN `seq` int NOT NULL AUTO_INCREMENT, ADD KEY `fk_seq` (`seq`)
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` ADD CONSTRAINT `fk_seq` FOREIGN KEY (`seq`) REFERENCES `fkmulti`.`p` (`id`)
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
|   `x` int(11) DEFAULT NULL,
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk_seq` FOREIGN KEY (`seq`) REFERENCES `fkmulti`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT ,
|   KEY `fk_seq` (`seq`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fkmulti_t_aig_fk_implicit_backing_drop
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_aig_fk_implicit_backing_drop DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_aig_fk_implicit_backing_drop
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT NOT NULL PRIMARY KEY, x INT)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_aig_fk_implicit_backing_drop' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_t_aig_fk_implicit_backing_drop.c
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti_t_aig_fk_implicit_backing_drop.p
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
< ok 3 0
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
| CREATE TABLE c (id INT NOT NULL PRIMARY KEY, x INT, seq INT NOT NULL AUTO_INCREMENT, CONSTRAINT fk_seq FOREIGN KEY (seq) REFERENCES p (id))
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
|   `x` int(11) DEFAULT NULL,
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk_seq` FOREIGN KEY (`seq`) REFERENCES `fkmulti`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
| ALTER TABLE `fkmulti`.`c` DROP FOREIGN KEY `fk_seq`
< ok 0 0
> exec
| ALTER TABLE `fkmulti`.`c` DROP INDEX `fk_seq`, DROP COLUMN `seq`
< error
| 1091 42000 Can't DROP 'fk_seq'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS fkmulti_t_aig_fk_rename_in_place_to_ai
< ok 0 0
> exec
| CREATE DATABASE fkmulti_t_aig_fk_rename_in_place_to_ai DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fkmulti_t_aig_fk_rename_in_place_to_ai
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT NOT NULL PRIMARY KEY, old_id INT, seq INT NOT NULL AUTO_INCREMENT, CONSTRAINT fk FOREIGN KEY (seq) REFERENCES p (id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'fkmulti_t_aig_fk_rename_in_place_to_ai' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| c
< row
: bytes
| p
> query
| SHOW CREATE TABLE fkmulti_t_aig_fk_rename_in_place_to_ai.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL,
|   `old_id` int(11) DEFAULT NULL,
|   `seq` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk` FOREIGN KEY (`seq`) REFERENCES `fkmulti_t_aig_fk_rename_in_place_to_ai`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fkmulti_t_aig_fk_rename_in_place_to_ai.p
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
| CREATE TABLE c (id INT NOT NULL PRIMARY KEY, old_id INT, CONSTRAINT fk FOREIGN KEY (old_id) REFERENCES p (id))
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
|   `old_id` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk` FOREIGN KEY (`old_id`) REFERENCES `fkmulti`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
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
| ALTER TABLE `fkmulti`.`c` DROP INDEX `fk`
< error
| 1091 42000 Can't DROP 'fk'; check that column/key exists
