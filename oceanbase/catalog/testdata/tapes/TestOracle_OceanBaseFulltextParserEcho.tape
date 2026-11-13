-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseFulltextParserEcho default#1
> exec
| DROP DATABASE IF EXISTS obft_0
< ok 0 0
> exec
| CREATE DATABASE obft_0
< ok 1 0
> exec
| USE obft_0
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c))
< ok 0 0
> query
| SHOW CREATE TABLE obft_0.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_0
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_1
< ok 0 0
> exec
| CREATE DATABASE obft_1
< ok 1 0
> exec
| USE obft_1
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER ngram)
< ok 0 0
> query
| SHOW CREATE TABLE obft_1.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=2) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_1
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_2
< ok 0 0
> exec
| CREATE DATABASE obft_2
< ok 1 0
> exec
| USE obft_2
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER ngram2)
< ok 0 0
> query
| SHOW CREATE TABLE obft_2.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER ngram2 PARSER_PROPERTIES=(min_ngram_size=2,max_ngram_size=3) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_2
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_3
< ok 0 0
> exec
| CREATE DATABASE obft_3
< ok 1 0
> exec
| USE obft_3
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER beng)
< ok 0 0
> query
| SHOW CREATE TABLE obft_3.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER beng PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_3
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_4
< ok 0 0
> exec
| CREATE DATABASE obft_4
< ok 1 0
> exec
| USE obft_4
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER ik)
< ok 0 0
> query
| SHOW CREATE TABLE obft_4.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER ik PARSER_PROPERTIES=(ik_mode="smart") BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_4
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_5
< ok 0 0
> exec
| CREATE DATABASE obft_5
< ok 1 0
> exec
| USE obft_5
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER space PARSER_PROPERTIES=(min_token_size=5))
< ok 0 0
> query
| SHOW CREATE TABLE obft_5.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=5,max_token_size=84) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_5
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_6
< ok 0 0
> exec
| CREATE DATABASE obft_6
< ok 1 0
> exec
| USE obft_6
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(max_token_size=50, min_token_size=2))
< ok 0 0
> query
| SHOW CREATE TABLE obft_6.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=2,max_token_size=50) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_6
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_7
< ok 0 0
> exec
| CREATE DATABASE obft_7
< ok 1 0
> exec
| USE obft_7
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3))
< ok 0 0
> query
| SHOW CREATE TABLE obft_7.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_7
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_8
< ok 0 0
> exec
| CREATE DATABASE obft_8
< ok 1 0
> exec
| USE obft_8
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) COMMENT 'x' WITH PARSER ngram LOCAL)
< ok 0 0
> query
| SHOW CREATE TABLE obft_8.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=2) BLOCK_SIZE 16384 COMMENT 'x'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_8
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_9
< ok 0 0
> exec
| CREATE DATABASE obft_9
< ok 1 0
> exec
| USE obft_9
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER NGRAM)
< ok 0 0
> query
| SHOW CREATE TABLE obft_9.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=2) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_9
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_10
< ok 0 0
> exec
| CREATE DATABASE obft_10
< ok 1 0
> exec
| USE obft_10
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER ik PARSER_PROPERTIES=(ik_mode='MAX_WORD'))
< ok 0 0
> query
| SHOW CREATE TABLE obft_10.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER ik PARSER_PROPERTIES=(ik_mode="max_word") BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_10
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_11
< ok 0 0
> exec
| CREATE DATABASE obft_11
< ok 1 0
> exec
| USE obft_11
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER ik PARSER_PROPERTIES=(ik_mode='max_word', stopword_table='none'))
< ok 0 0
> query
| SHOW CREATE TABLE obft_11.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER ik PARSER_PROPERTIES=(ik_mode="max_word") BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_11
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_12
< ok 0 0
> exec
| CREATE DATABASE obft_12
< ok 1 0
> exec
| USE obft_12
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER ngram2 PARSER_PROPERTIES=(min_ngram_size=1, max_ngram_size=4))
< ok 0 0
> query
| SHOW CREATE TABLE obft_12.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER ngram2 PARSER_PROPERTIES=(min_ngram_size=1,max_ngram_size=4) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_12
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_13
< ok 0 0
> exec
| CREATE DATABASE obft_13
< ok 1 0
> exec
| USE obft_13
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER beng PARSER_PROPERTIES=(min_token_size=4,max_token_size=20) INVISIBLE)
< ok 0 0
> query
| SHOW CREATE TABLE obft_13.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER beng PARSER_PROPERTIES=(min_token_size=4,max_token_size=20) BLOCK_SIZE 16384 /*!80000 INVISIBLE */
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_13
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_14
< ok 0 0
> exec
| CREATE DATABASE obft_14
< ok 1 0
> exec
| USE obft_14
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c, d) WITH PARSER space PARSER_PROPERTIES (min_token_size=2 max_token_size=80))
< ok 0 0
> query
| SHOW CREATE TABLE obft_14.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`, `d`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=2,max_token_size=80) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_14
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_15
< ok 0 0
> exec
| CREATE DATABASE obft_15
< ok 1 0
> exec
| USE obft_15
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(min_token_size=5) PARSER_PROPERTIES=(max_token_size=40))
< ok 0 0
> query
| SHOW CREATE TABLE obft_15.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=40) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_15
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_16
< ok 0 0
> exec
| CREATE DATABASE obft_16
< ok 1 0
> exec
| USE obft_16
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3, min_token_size=6))
< ok 0 0
> query
| SHOW CREATE TABLE obft_16.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=6,max_token_size=84) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_16
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_17
< ok 0 0
> exec
| CREATE DATABASE obft_17
< ok 1 0
> exec
| USE obft_17
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER ngram WITH PARSER space)
< ok 0 0
> query
| SHOW CREATE TABLE obft_17.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_17
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_18
< ok 0 0
> exec
| CREATE DATABASE obft_18
< ok 1 0
> exec
| USE obft_18
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(ngram_token_size=3) WITH PARSER ngram)
< ok 0 0
> query
| SHOW CREATE TABLE obft_18.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_18
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obft_like
< ok 0 0
> exec
| CREATE DATABASE obft_like
< ok 1 0
> exec
| USE obft_like
< ok 0 0
> exec
| CREATE TABLE src (id INT NOT NULL, c VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3))
< ok 0 0
> query
| SHOW CREATE TABLE obft_like.src
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| src
: bytes
| CREATE TABLE `src` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE TABLE t LIKE src
< ok 0 0
> query
| SHOW CREATE TABLE obft_like.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obft_like
< ok 10 0
> exec
| DROP DATABASE IF EXISTS obftchg_0
< ok 0 0
> exec
| CREATE DATABASE obftchg_0
< ok 1 0
> exec
| USE obftchg_0
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c))
< ok 0 0
> query
| SHOW CREATE TABLE obftchg_0.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obftchg_0
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obftchg_1
< ok 0 0
> exec
| CREATE DATABASE obftchg_1
< ok 1 0
> exec
| USE obftchg_1
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER ngram)
< ok 0 0
> query
| SHOW CREATE TABLE obftchg_1.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=2) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obftchg_1
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obftchg_2
< ok 0 0
> exec
| CREATE DATABASE obftchg_2
< ok 1 0
> exec
| USE obftchg_2
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER ik)
< ok 0 0
> query
| SHOW CREATE TABLE obftchg_2.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER ik PARSER_PROPERTIES=(ik_mode="smart") BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obftchg_2
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obftchg_3
< ok 0 0
> exec
| CREATE DATABASE obftchg_3
< ok 1 0
> exec
| USE obftchg_3
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c))
< ok 0 0
> query
| SHOW CREATE TABLE obftchg_3.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obftchg_3
< ok 5 0
> exec
| DROP DATABASE IF EXISTS obftchg_4
< ok 0 0
> exec
| CREATE DATABASE obftchg_4
< ok 1 0
> exec
| USE obftchg_4
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c VARCHAR(100), d VARCHAR(100), PRIMARY KEY (id), FULLTEXT KEY fk (c) WITH PARSER space)
< ok 0 0
> query
| SHOW CREATE TABLE obftchg_4.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `c` varchar(100) DEFAULT NULL,
|   `d` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `fk` (`c`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obftchg_4
< ok 5 0
