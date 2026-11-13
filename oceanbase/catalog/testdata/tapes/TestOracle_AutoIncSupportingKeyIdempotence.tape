-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_AutoIncSupportingKeyIdempotence default#1
> exec
| DROP DATABASE IF EXISTS aidem_autoinc_nonpk_unique
< ok 0 0
> exec
| CREATE DATABASE aidem_autoinc_nonpk_unique
< ok 1 0
> exec
| USE aidem_autoinc_nonpk_unique
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, name VARCHAR(20), UNIQUE KEY uk (id)) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE aidem_autoinc_nonpk_unique.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `name` varchar(20) DEFAULT NULL,
|   UNIQUE KEY `uk` (`id`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS aidem_autoinc_nonpk_unique_composite
< ok 0 0
> exec
| CREATE DATABASE aidem_autoinc_nonpk_unique_composite
< ok 1 0
> exec
| USE aidem_autoinc_nonpk_unique_composite
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk (id, a)) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE aidem_autoinc_nonpk_unique_composite.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `a` int(11) NOT NULL,
|   UNIQUE KEY `uk` (`id`, `a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS aidem_autoinc_plain_key
< ok 0 0
> exec
| CREATE DATABASE aidem_autoinc_plain_key
< ok 1 0
> exec
| USE aidem_autoinc_plain_key
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT, KEY k (id)) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE aidem_autoinc_plain_key.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `a` int(11) DEFAULT NULL,
|   KEY `k` (`id`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS aidem_autoinc_unique_with_other_pk
< ok 0 0
> exec
| CREATE DATABASE aidem_autoinc_unique_with_other_pk
< ok 1 0
> exec
| USE aidem_autoinc_unique_with_other_pk
< ok 0 0
> exec
| CREATE TABLE t (pk INT NOT NULL, id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (pk), UNIQUE KEY uk (id)) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE aidem_autoinc_unique_with_other_pk.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `pk` int(11) NOT NULL,
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`pk`),
|   UNIQUE KEY `uk` (`id`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS aidem_autoinc_primary_key_guard
< ok 0 0
> exec
| CREATE DATABASE aidem_autoinc_primary_key_guard
< ok 1 0
> exec
| USE aidem_autoinc_primary_key_guard
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, name VARCHAR(20)) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE aidem_autoinc_primary_key_guard.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
