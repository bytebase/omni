-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_IndexMigrationApplyCorrectness default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_secondary_3dc43f1e_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_secondary_3dc43f1e_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_secondary_3dc43f1e_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ka (a), KEY kab (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_secondary_3dc43f1e_to.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `kab` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_secondary_3dc43f1e
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_secondary_3dc43f1e DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_secondary_3dc43f1e
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_secondary_3dc43f1e`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_secondary_3dc43f1e`.`t` ADD KEY `ka` (`a`) LOCAL
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_secondary_3dc43f1e`.`t` ADD KEY `kab` (`a`,`b`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_secondary_3dc43f1e.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `kab` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_unique_39a3093b_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_unique_39a3093b_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_unique_39a3093b_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, UNIQUE KEY ua (a), UNIQUE KEY uab (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_unique_39a3093b_to.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `uab` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_unique_39a3093b
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_unique_39a3093b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_unique_39a3093b
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_unique_39a3093b`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_unique_39a3093b`.`t` ADD UNIQUE KEY `ua` (`a`) LOCAL
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_unique_39a3093b`.`t` ADD UNIQUE KEY `uab` (`a`,`b`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_unique_39a3093b.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `uab` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_prefix_5a0adfb0_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_prefix_5a0adfb0_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_prefix_5a0adfb0_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100), KEY sp (s(10)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_prefix_5a0adfb0_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `sp` (`s`(10)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_prefix_5a0adfb0
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_prefix_5a0adfb0 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_prefix_5a0adfb0
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_prefix_5a0adfb0`.`t` (
|   `id` int NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_prefix_5a0adfb0`.`t` ADD KEY `sp` (`s`(10)) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_prefix_5a0adfb0.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `sp` (`s`(10)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_comment_a6dadac1_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_comment_a6dadac1_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_comment_a6dadac1_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a) COMMENT 'hello')
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_comment_a6dadac1_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL COMMENT 'hello'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_comment_a6dadac1
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_comment_a6dadac1 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_comment_a6dadac1
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_comment_a6dadac1`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_comment_a6dadac1`.`t` ADD KEY `ka` (`a`) COMMENT 'hello' LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_comment_a6dadac1.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL COMMENT 'hello'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_using_btree_4c930a0b_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_using_btree_4c930a0b_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_using_btree_4c930a0b_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ua (a) USING BTREE) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_using_btree_4c930a0b_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `ua` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_using_btree_4c930a0b
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_using_btree_4c930a0b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_using_btree_4c930a0b
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_using_btree_4c930a0b`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_using_btree_4c930a0b`.`t` ADD KEY `ua` (`a`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_using_btree_4c930a0b.t
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
|   PRIMARY KEY (`id`),
|   KEY `ua` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_kbs_a4a590ca_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_kbs_a4a590ca_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_kbs_a4a590ca_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a) KEY_BLOCK_SIZE=4) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_kbs_a4a590ca_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_kbs_a4a590ca
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_kbs_a4a590ca DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_kbs_a4a590ca
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_kbs_a4a590ca`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_kbs_a4a590ca`.`t` ADD KEY `ka` (`a`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_kbs_a4a590ca.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_composite_pk_d5b7ea43_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_composite_pk_d5b7ea43_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_composite_pk_d5b7ea43_to
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, c INT, PRIMARY KEY (a,b), KEY kc (c))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_composite_pk_d5b7ea43_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   PRIMARY KEY (`a`, `b`),
|   KEY `kc` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_composite_pk_d5b7ea43
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_composite_pk_d5b7ea43 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_composite_pk_d5b7ea43
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_composite_pk_d5b7ea43`.`t` (
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_composite_pk_d5b7ea43`.`t` ADD KEY `kc` (`c`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_composite_pk_d5b7ea43.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   PRIMARY KEY (`a`, `b`),
|   KEY `kc` (`c`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_fulltext_1a5e79fa_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_fulltext_1a5e79fa_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_fulltext_1a5e79fa_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, body TEXT, FULLTEXT KEY ftb (body))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_fulltext_1a5e79fa_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `body` text DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `ftb` (`body`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_fulltext_1a5e79fa
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_fulltext_1a5e79fa DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_fulltext_1a5e79fa
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_fulltext_1a5e79fa`.`t` (
|   `id` int NOT NULL,
|   `body` text DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_fulltext_1a5e79fa`.`t` ADD FULLTEXT KEY `ftb` (`body`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_fulltext_1a5e79fa.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `body` text DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `ftb` (`body`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_descending_c77de006_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_descending_c77de006_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_descending_c77de006_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ab (a, b DESC))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_descending_c77de006_to.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ab` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_descending_c77de006
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_descending_c77de006 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_descending_c77de006
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_descending_c77de006`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_descending_c77de006`.`t` ADD KEY `ab` (`a`,`b`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_descending_c77de006.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ab` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_invisible_caa5ab2f_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_invisible_caa5ab2f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_invisible_caa5ab2f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY va (a) /*!80000 INVISIBLE */)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_invisible_caa5ab2f_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `va` (`a`) BLOCK_SIZE 16384 LOCAL /*!80000 INVISIBLE */
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_invisible_caa5ab2f
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_invisible_caa5ab2f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_invisible_caa5ab2f
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_invisible_caa5ab2f`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_invisible_caa5ab2f`.`t` ADD KEY `va` (`a`) /*!80000 INVISIBLE */ LOCAL
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_invisible_caa5ab2f`.`t` ALTER INDEX `va` INVISIBLE
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_invisible_caa5ab2f.t
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
|   PRIMARY KEY (`id`),
|   KEY `va` (`a`) BLOCK_SIZE 16384 LOCAL /*!80000 INVISIBLE */
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_functional_1f811dd1_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_functional_1f811dd1_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_functional_1f811dd1_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY fa ((a + 1)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_functional_1f811dd1_to.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fa` ((`a` + 1)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_functional_1f811dd1
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_functional_1f811dd1 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_functional_1f811dd1
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_functional_1f811dd1`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_functional_1f811dd1`.`t` ADD KEY `fa` (((`a` + 1))) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_functional_1f811dd1.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fa` ((`a` + 1)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_mv_index_8238bf32_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_mv_index_8238bf32_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_mv_index_8238bf32_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, j JSON, KEY mv ((CAST(j->'$.ids' AS UNSIGNED ARRAY))))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_mv_index_8238bf32_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   INDEX `mv` ((json_query(`j`,'$.ids' RETURNING unsigned WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch))) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_mv_index_8238bf32
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_mv_index_8238bf32 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_mv_index_8238bf32
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_mv_index_8238bf32`.`t` (
|   `id` int NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `mv` ((json_query(`j`,_utf8mb4'$.ids' RETURNING unsigned WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_mv_index_8238bf32.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   INDEX `mv` ((json_query(`j`,'$.ids' RETURNING unsigned WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch))) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_spatial_6826187c_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_spatial_6826187c_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_spatial_6826187c_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL, SPATIAL KEY sg (g))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_spatial_6826187c_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `sg` (`g`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_spatial_6826187c
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_spatial_6826187c DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_spatial_6826187c
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_spatial_6826187c`.`t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_spatial_6826187c`.`t` ADD SPATIAL KEY `sg` (`g`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_spatial_6826187c.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `sg` (`g`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_secondary_5a1931ba_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_secondary_5a1931ba_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_secondary_5a1931ba_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_secondary_5a1931ba_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_secondary_5a1931ba
< ok 0 0
> exec
| CREATE DATABASE pa_add_secondary_5a1931ba DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_secondary_5a1931ba
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_secondary_5a1931ba.t
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
| ALTER TABLE `pa_add_secondary_5a1931ba`.`t` ADD KEY `ka` (`a`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_secondary_5a1931ba.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_unique_c33bac8f_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_unique_c33bac8f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_unique_c33bac8f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, UNIQUE KEY ua (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_unique_c33bac8f_to.t
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
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_unique_c33bac8f
< ok 0 0
> exec
| CREATE DATABASE pa_add_unique_c33bac8f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_unique_c33bac8f
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_unique_c33bac8f.t
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
| ALTER TABLE `pa_add_unique_c33bac8f`.`t` ADD UNIQUE KEY `ua` (`a`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_unique_c33bac8f.t
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
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_composite_30ce947b_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_composite_30ce947b_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_composite_30ce947b_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY kab (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_composite_30ce947b_to.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `kab` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_composite_30ce947b
< ok 0 0
> exec
| CREATE DATABASE pa_add_composite_30ce947b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_composite_30ce947b
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_composite_30ce947b.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_add_composite_30ce947b`.`t` ADD KEY `kab` (`a`,`b`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_composite_30ce947b.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `kab` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_prefix_d19cf194_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_prefix_d19cf194_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_prefix_d19cf194_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100), KEY sp (s(10)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_prefix_d19cf194_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `sp` (`s`(10)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_prefix_d19cf194
< ok 0 0
> exec
| CREATE DATABASE pa_add_prefix_d19cf194 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_prefix_d19cf194
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_prefix_d19cf194.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_add_prefix_d19cf194`.`t` ADD KEY `sp` (`s`(10)) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_prefix_d19cf194.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `sp` (`s`(10)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_comment_3f58f4ad_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_comment_3f58f4ad_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_comment_3f58f4ad_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a) COMMENT 'c')
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_comment_3f58f4ad_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL COMMENT 'c'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_comment_3f58f4ad
< ok 0 0
> exec
| CREATE DATABASE pa_add_comment_3f58f4ad DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_comment_3f58f4ad
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_comment_3f58f4ad.t
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
| ALTER TABLE `pa_add_comment_3f58f4ad`.`t` ADD KEY `ka` (`a`) COMMENT 'c' LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_comment_3f58f4ad.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL COMMENT 'c'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_fulltext_151ea6a6_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_fulltext_151ea6a6_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_fulltext_151ea6a6_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, body TEXT, FULLTEXT KEY ftb (body))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_fulltext_151ea6a6_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `body` text DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `ftb` (`body`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_fulltext_151ea6a6
< ok 0 0
> exec
| CREATE DATABASE pa_add_fulltext_151ea6a6 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_fulltext_151ea6a6
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, body TEXT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_fulltext_151ea6a6.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `body` text DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_add_fulltext_151ea6a6`.`t` ADD FULLTEXT KEY `ftb` (`body`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_fulltext_151ea6a6.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `body` text DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `ftb` (`body`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_primary_key_f0223c58_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_primary_key_f0223c58_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_primary_key_f0223c58_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_primary_key_f0223c58_to.t
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
| DROP DATABASE IF EXISTS pa_add_primary_key_f0223c58
< ok 0 0
> exec
| CREATE DATABASE pa_add_primary_key_f0223c58 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_primary_key_f0223c58
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_primary_key_f0223c58.t
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
| ALTER TABLE `pa_add_primary_key_f0223c58`.`t` ADD PRIMARY KEY (`id`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_primary_key_f0223c58.t
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
| DROP DATABASE IF EXISTS pa_add_descending_b60eb382_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_descending_b60eb382_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_descending_b60eb382_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ab (a, b DESC))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_descending_b60eb382_to.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ab` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_descending_b60eb382
< ok 0 0
> exec
| CREATE DATABASE pa_add_descending_b60eb382 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_descending_b60eb382
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_descending_b60eb382.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_add_descending_b60eb382`.`t` ADD KEY `ab` (`a`,`b`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_descending_b60eb382.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ab` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_invisible_5137783b_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_invisible_5137783b_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_invisible_5137783b_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY va (a) /*!80000 INVISIBLE */)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_invisible_5137783b_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `va` (`a`) BLOCK_SIZE 16384 LOCAL /*!80000 INVISIBLE */
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_invisible_5137783b
< ok 0 0
> exec
| CREATE DATABASE pa_add_invisible_5137783b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_invisible_5137783b
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_invisible_5137783b.t
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
| ALTER TABLE `pa_add_invisible_5137783b`.`t` ADD KEY `va` (`a`) /*!80000 INVISIBLE */ LOCAL
< ok 0 0
> exec
| ALTER TABLE `pa_add_invisible_5137783b`.`t` ALTER INDEX `va` INVISIBLE
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_invisible_5137783b.t
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
|   PRIMARY KEY (`id`),
|   KEY `va` (`a`) BLOCK_SIZE 16384 LOCAL /*!80000 INVISIBLE */
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_functional_fccddf0d_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_functional_fccddf0d_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_functional_fccddf0d_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY fa ((a + 1)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_functional_fccddf0d_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `fa` ((`a` + 1)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_functional_fccddf0d
< ok 0 0
> exec
| CREATE DATABASE pa_add_functional_fccddf0d DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_functional_fccddf0d
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_functional_fccddf0d.t
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
| ALTER TABLE `pa_add_functional_fccddf0d`.`t` ADD KEY `fa` (((`a` + 1))) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_functional_fccddf0d.t
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
|   PRIMARY KEY (`id`),
|   KEY `fa` ((`a` + 1)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_mv_index_dbbc6d6_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_mv_index_dbbc6d6_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_mv_index_dbbc6d6_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, j JSON, KEY mv ((CAST(j->'$.ids' AS UNSIGNED ARRAY))))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_mv_index_dbbc6d6_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   INDEX `mv` ((json_query(`j`,'$.ids' RETURNING unsigned WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch))) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_mv_index_dbbc6d6
< ok 0 0
> exec
| CREATE DATABASE pa_add_mv_index_dbbc6d6 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_mv_index_dbbc6d6
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, j JSON)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_mv_index_dbbc6d6.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_add_mv_index_dbbc6d6`.`t` ADD KEY `mv` ((json_query(`j`,_utf8mb4'$.ids' RETURNING unsigned WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch)))
< error
| 1235 0A000 build multivalue index afterward not supported
> exec
| DROP DATABASE IF EXISTS pa_drop_secondary_a80bd9dc_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_secondary_a80bd9dc_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_secondary_a80bd9dc_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_secondary_a80bd9dc_to.t
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
| DROP DATABASE IF EXISTS pa_drop_secondary_a80bd9dc
< ok 0 0
> exec
| CREATE DATABASE pa_drop_secondary_a80bd9dc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_secondary_a80bd9dc
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_secondary_a80bd9dc.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_drop_secondary_a80bd9dc`.`t` DROP INDEX `ka`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_secondary_a80bd9dc.t
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
| DROP DATABASE IF EXISTS pa_drop_unique_17dca13d_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_unique_17dca13d_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_unique_17dca13d_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_unique_17dca13d_to.t
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
| DROP DATABASE IF EXISTS pa_drop_unique_17dca13d
< ok 0 0
> exec
| CREATE DATABASE pa_drop_unique_17dca13d DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_unique_17dca13d
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, UNIQUE KEY ua (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_unique_17dca13d.t
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
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_drop_unique_17dca13d`.`t` DROP INDEX `ua`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_unique_17dca13d.t
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
| DROP DATABASE IF EXISTS pa_drop_primary_key_8793bdb2_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_primary_key_8793bdb2_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_primary_key_8793bdb2_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_primary_key_8793bdb2_to.t
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
| DROP DATABASE IF EXISTS pa_drop_primary_key_8793bdb2
< ok 0 0
> exec
| CREATE DATABASE pa_drop_primary_key_8793bdb2 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_primary_key_8793bdb2
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_primary_key_8793bdb2.t
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
| ALTER TABLE `pa_drop_primary_key_8793bdb2`.`t` DROP PRIMARY KEY
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_primary_key_8793bdb2.t
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
| DROP DATABASE IF EXISTS pa_drop_one_of_many_bd4ebc08_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_one_of_many_bd4ebc08_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_one_of_many_bd4ebc08_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY kb (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_one_of_many_bd4ebc08_to.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `kb` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_drop_one_of_many_bd4ebc08
< ok 0 0
> exec
| CREATE DATABASE pa_drop_one_of_many_bd4ebc08 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_one_of_many_bd4ebc08
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ka (a), KEY kb (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_one_of_many_bd4ebc08.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `kb` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_drop_one_of_many_bd4ebc08`.`t` DROP INDEX `ka`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_one_of_many_bd4ebc08.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `kb` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_drop_mv_index_410d3fe0_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_mv_index_410d3fe0_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_mv_index_410d3fe0_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, j JSON)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_mv_index_410d3fe0_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_drop_mv_index_410d3fe0
< ok 0 0
> exec
| CREATE DATABASE pa_drop_mv_index_410d3fe0 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_mv_index_410d3fe0
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, j JSON, KEY mv ((CAST(j->'$.ids' AS UNSIGNED ARRAY))))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_mv_index_410d3fe0.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   INDEX `mv` ((json_query(`j`,'$.ids' RETURNING unsigned WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch))) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_drop_mv_index_410d3fe0`.`t` DROP INDEX `mv`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_mv_index_410d3fe0.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_columns_92b80ce_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_columns_92b80ce_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_columns_92b80ce_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY k (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_columns_92b80ce_to.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_columns_92b80ce
< ok 0 0
> exec
| CREATE DATABASE pa_modify_columns_92b80ce DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_columns_92b80ce
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY k (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_columns_92b80ce.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_columns_92b80ce`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_columns_92b80ce`.`t` ADD KEY `k` (`a`,`b`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_columns_92b80ce.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_kind_to_unique_e8b5118f_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_kind_to_unique_e8b5118f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_kind_to_unique_e8b5118f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, UNIQUE KEY k (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_kind_to_unique_e8b5118f_to.t
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
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_kind_to_unique_e8b5118f
< ok 0 0
> exec
| CREATE DATABASE pa_modify_kind_to_unique_e8b5118f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_kind_to_unique_e8b5118f
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_kind_to_unique_e8b5118f.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_kind_to_unique_e8b5118f`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_kind_to_unique_e8b5118f`.`t` ADD UNIQUE KEY `k` (`a`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_kind_to_unique_e8b5118f.t
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
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_prefix_length_59d4710a_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_prefix_length_59d4710a_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_prefix_length_59d4710a_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100), KEY k (s(20)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_prefix_length_59d4710a_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`s`(20)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_prefix_length_59d4710a
< ok 0 0
> exec
| CREATE DATABASE pa_modify_prefix_length_59d4710a DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_prefix_length_59d4710a
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100), KEY k (s(10)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_prefix_length_59d4710a.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`s`(10)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_prefix_length_59d4710a`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_prefix_length_59d4710a`.`t` ADD KEY `k` (`s`(20)) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_prefix_length_59d4710a.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`s`(20)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_comment_a51be4e4_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_comment_a51be4e4_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_comment_a51be4e4_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) COMMENT 'new')
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_comment_a51be4e4_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL COMMENT 'new'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_comment_a51be4e4
< ok 0 0
> exec
| CREATE DATABASE pa_modify_comment_a51be4e4 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_comment_a51be4e4
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) COMMENT 'old')
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_comment_a51be4e4.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL COMMENT 'old'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_comment_a51be4e4`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_comment_a51be4e4`.`t` ADD KEY `k` (`a`) COMMENT 'new' LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_comment_a51be4e4.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL COMMENT 'new'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_columns_6203f0d8_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_columns_6203f0d8_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_columns_6203f0d8_to
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_columns_6203f0d8_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`a`, `b`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_columns_6203f0d8
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_columns_6203f0d8 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_columns_6203f0d8
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_columns_6203f0d8.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`a`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_pk_columns_6203f0d8`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`a`,`b`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_columns_6203f0d8.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`a`, `b`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_autoincrement_29ca8be3_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_autoincrement_29ca8be3_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_autoincrement_29ca8be3_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, PRIMARY KEY (id,a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_autoincrement_29ca8be3_to.t
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
|   PRIMARY KEY (`id`, `a`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_autoincrement_29ca8be3
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_autoincrement_29ca8be3 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_autoincrement_29ca8be3
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_autoincrement_29ca8be3.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_pk_autoincrement_29ca8be3`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`id`,`a`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_autoincrement_29ca8be3.t
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
|   PRIMARY KEY (`id`, `a`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_swap_autoincrement_b1e7f47f_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_swap_autoincrement_b1e7f47f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_swap_autoincrement_b1e7f47f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (id,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_swap_autoincrement_b1e7f47f_to.t
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
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`id`, `b`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_swap_autoincrement_b1e7f47f
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_swap_autoincrement_b1e7f47f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_swap_autoincrement_b1e7f47f
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (id,a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_swap_autoincrement_b1e7f47f.t
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
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`id`, `a`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_pk_swap_autoincrement_b1e7f47f`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`id`,`b`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_swap_autoincrement_b1e7f47f.t
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
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`id`, `b`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_drop_old_member_6f193c55_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_drop_old_member_6f193c55_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_drop_old_member_6f193c55_to
< ok 0 0
> exec
| CREATE TABLE t (b INT NOT NULL, c INT, PRIMARY KEY (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_drop_old_member_6f193c55_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_drop_old_member_6f193c55
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_drop_old_member_6f193c55 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_drop_old_member_6f193c55
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_drop_old_member_6f193c55.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`a`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_pk_drop_old_member_6f193c55`.`t` ADD COLUMN `c` int DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_modify_pk_drop_old_member_6f193c55`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`b`)
< ok 0 0
> exec
| ALTER TABLE `pa_modify_pk_drop_old_member_6f193c55`.`t` DROP COLUMN `a`
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_drop_old_member_6f193c55.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_shrink_demote_nullable_c7332ccf_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_shrink_demote_nullable_c7332ccf_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_shrink_demote_nullable_c7332ccf_to
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NULL, PRIMARY KEY (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_shrink_demote_nullable_c7332ccf_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`a`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_pk_shrink_demote_nullable_c7332ccf
< ok 0 0
> exec
| CREATE DATABASE pa_modify_pk_shrink_demote_nullable_c7332ccf DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_pk_shrink_demote_nullable_c7332ccf
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_shrink_demote_nullable_c7332ccf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`a`, `b`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_pk_shrink_demote_nullable_c7332ccf`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`a`)
< ok 0 0
> exec
| ALTER TABLE `pa_modify_pk_shrink_demote_nullable_c7332ccf`.`t` MODIFY COLUMN `b` int DEFAULT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_pk_shrink_demote_nullable_c7332ccf.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`a`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_visibility_b465b887_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_visibility_b465b887_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_visibility_b465b887_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) /*!80000 INVISIBLE */)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_visibility_b465b887_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL /*!80000 INVISIBLE */
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_visibility_b465b887
< ok 0 0
> exec
| CREATE DATABASE pa_modify_visibility_b465b887 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_visibility_b465b887
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_visibility_b465b887.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_visibility_b465b887`.`t` ALTER INDEX `k` INVISIBLE
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_visibility_b465b887.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL /*!80000 INVISIBLE */
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_direction_c8075b6c_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_direction_c8075b6c_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_direction_c8075b6c_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a DESC))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_direction_c8075b6c_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_direction_c8075b6c
< ok 0 0
> exec
| CREATE DATABASE pa_modify_direction_c8075b6c DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_direction_c8075b6c
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_direction_c8075b6c.t
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query ^
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
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_drop_column_and_its_single_col_index_afb3e4c0_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_column_and_its_single_col_index_afb3e4c0_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_column_and_its_single_col_index_afb3e4c0_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_column_and_its_single_col_index_afb3e4c0_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_drop_column_and_its_single_col_index_afb3e4c0
< ok 0 0
> exec
| CREATE DATABASE pa_drop_column_and_its_single_col_index_afb3e4c0 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_column_and_its_single_col_index_afb3e4c0
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_column_and_its_single_col_index_afb3e4c0.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_drop_column_and_its_single_col_index_afb3e4c0`.`t` DROP INDEX `ka`
< ok 0 0
> exec
| ALTER TABLE `pa_drop_column_and_its_single_col_index_afb3e4c0`.`t` DROP COLUMN `a`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_column_and_its_single_col_index_afb3e4c0.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_drop_columns_and_their_composite_index_31114bc7_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_columns_and_their_composite_index_31114bc7_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_columns_and_their_composite_index_31114bc7_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_columns_and_their_composite_index_31114bc7_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_drop_columns_and_their_composite_index_31114bc7
< ok 0 0
> exec
| CREATE DATABASE pa_drop_columns_and_their_composite_index_31114bc7 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_columns_and_their_composite_index_31114bc7
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY kab (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_columns_and_their_composite_index_31114bc7.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `kab` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_drop_columns_and_their_composite_index_31114bc7`.`t` DROP INDEX `kab`
< ok 0 0
> exec
| ALTER TABLE `pa_drop_columns_and_their_composite_index_31114bc7`.`t` DROP COLUMN `a`
< ok 0 0
> exec
| ALTER TABLE `pa_drop_columns_and_their_composite_index_31114bc7`.`t` DROP COLUMN `b`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_columns_and_their_composite_index_31114bc7.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_index_shrink_with_column_drop_41aa6ae5_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_index_shrink_with_column_drop_41aa6ae5_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_index_shrink_with_column_drop_41aa6ae5_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, b INT, KEY k (b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_index_shrink_with_column_drop_41aa6ae5_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_index_shrink_with_column_drop_41aa6ae5
< ok 0 0
> exec
| CREATE DATABASE pa_modify_index_shrink_with_column_drop_41aa6ae5 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_index_shrink_with_column_drop_41aa6ae5
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY k (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_index_shrink_with_column_drop_41aa6ae5.t
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
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_index_shrink_with_column_drop_41aa6ae5`.`t` DROP INDEX `k`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_index_shrink_with_column_drop_41aa6ae5`.`t` DROP COLUMN `a`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_index_shrink_with_column_drop_41aa6ae5`.`t` ADD KEY `k` (`b`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_index_shrink_with_column_drop_41aa6ae5.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_column_and_index_on_it_dcc5d6b_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_and_index_on_it_dcc5d6b_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_and_index_on_it_dcc5d6b_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_and_index_on_it_dcc5d6b_to.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_column_and_index_on_it_dcc5d6b
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_and_index_on_it_dcc5d6b DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_and_index_on_it_dcc5d6b
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_and_index_on_it_dcc5d6b.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_add_column_and_index_on_it_dcc5d6b`.`t` ADD COLUMN `a` int DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_add_column_and_index_on_it_dcc5d6b`.`t` ADD KEY `ka` (`a`) LOCAL
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_and_index_on_it_dcc5d6b.t
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
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
