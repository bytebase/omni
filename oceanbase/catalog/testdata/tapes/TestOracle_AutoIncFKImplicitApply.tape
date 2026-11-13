-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_AutoIncFKImplicitApply default#1
> exec
| DROP DATABASE IF EXISTS aifk_autoinc_fk_implicit_only
< ok 0 0
> exec
| CREATE DATABASE aifk_autoinc_fk_implicit_only DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE aifk_autoinc_fk_implicit_only
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (id INT NOT NULL AUTO_INCREMENT, x INT, FOREIGN KEY (id) REFERENCES p(id)) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE aifk_autoinc_fk_implicit_only.p
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
| SHOW CREATE TABLE aifk_autoinc_fk_implicit_only.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `x` int(11) DEFAULT NULL,
|   CONSTRAINT `c_OBFK_0000000000000001` FOREIGN KEY (`id`) REFERENCES `aifk_autoinc_fk_implicit_only`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS aifk_autoinc_fk_implicit_only
< ok 2 0
> exec
| CREATE DATABASE aifk_autoinc_fk_implicit_only DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE aifk_autoinc_fk_implicit_only
< ok 0 0
> exec
| CREATE TABLE `aifk_autoinc_fk_implicit_only`.`c` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `x` int DEFAULT NULL,
|   KEY `c_OBFK_0000000000000001` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| CREATE TABLE `aifk_autoinc_fk_implicit_only`.`p` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `aifk_autoinc_fk_implicit_only`.`c` ADD CONSTRAINT `c_OBFK_0000000000000001` FOREIGN KEY (`id`) REFERENCES `aifk_autoinc_fk_implicit_only`.`p` (`id`)
< ok 0 0
> query
| SHOW CREATE TABLE aifk_autoinc_fk_implicit_only.p
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
| SHOW CREATE TABLE aifk_autoinc_fk_implicit_only.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `x` int(11) DEFAULT NULL,
|   CONSTRAINT `c_OBFK_0000000000000001` FOREIGN KEY (`id`) REFERENCES `aifk_autoinc_fk_implicit_only`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT ,
|   KEY `c_OBFK_0000000000000001` (`id`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
