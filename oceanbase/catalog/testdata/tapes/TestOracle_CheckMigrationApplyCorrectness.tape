-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_CheckMigrationApplyCorrectness default#1
> exec
| DROP DATABASE IF EXISTS pa_create_with_named_check_93d8220_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_named_check_93d8220_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_named_check_93d8220_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_named_check_93d8220_to.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_named_check_93d8220
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_named_check_93d8220 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_named_check_93d8220
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_named_check_93d8220`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_named_check_93d8220`.`t` ADD CONSTRAINT `ck` CHECK ((`a` > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_named_check_93d8220.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_unnamed_check_995eb16d_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_unnamed_check_995eb16d_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_unnamed_check_995eb16d_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_unnamed_check_995eb16d_to.t
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
|   CONSTRAINT `t_OBCHECK_0000000000000001` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_unnamed_check_995eb16d
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_unnamed_check_995eb16d DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_unnamed_check_995eb16d
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_unnamed_check_995eb16d`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_unnamed_check_995eb16d`.`t` ADD CHECK ((`a` > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_unnamed_check_995eb16d.t
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
|   CONSTRAINT `t_OBCHECK_0000000000000002` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_multiple_checks_a86c99ea_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_multiple_checks_a86c99ea_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_multiple_checks_a86c99ea_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a > 0), CONSTRAINT cb CHECK (b < 100), CHECK (a < b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_multiple_checks_a86c99ea_to.t
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
|   CONSTRAINT `t_OBCHECK_0000000000000003` CHECK ((`a` > 0)),
|   CONSTRAINT `cb` CHECK ((`b` < 100)),
|   CONSTRAINT `t_OBCHECK_0000000000000004` CHECK ((`a` < `b`))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_multiple_checks_a86c99ea
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_multiple_checks_a86c99ea DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_multiple_checks_a86c99ea
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_multiple_checks_a86c99ea`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_multiple_checks_a86c99ea`.`t` ADD CONSTRAINT `cb` CHECK ((`b` < 100))
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_multiple_checks_a86c99ea`.`t` ADD CHECK ((`a` > 0))
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_multiple_checks_a86c99ea`.`t` ADD CHECK ((`a` < `b`))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_multiple_checks_a86c99ea.t
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
|   CONSTRAINT `cb` CHECK ((`b` < 100)),
|   CONSTRAINT `t_OBCHECK_0000000000000005` CHECK ((`a` > 0)),
|   CONSTRAINT `t_OBCHECK_0000000000000006` CHECK ((`a` < `b`))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_not_enforced_79c7c4f2_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_not_enforced_79c7c4f2_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_not_enforced_79c7c4f2_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ne CHECK (a < b) NOT ENFORCED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_not_enforced_79c7c4f2_to.t
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
|   CONSTRAINT `ne` CHECK ((`a` < `b`))  NOT ENFORCED
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_not_enforced_79c7c4f2
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_not_enforced_79c7c4f2 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_not_enforced_79c7c4f2
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_not_enforced_79c7c4f2`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_not_enforced_79c7c4f2`.`t` ADD CONSTRAINT `ne` CHECK ((`a` < `b`)) /*!80016 NOT ENFORCED */
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_not_enforced_79c7c4f2.t
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
|   CONSTRAINT `ne` CHECK ((`a` < `b`))  NOT ENFORCED
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_compound_40e1387d_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_compound_40e1387d_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_compound_40e1387d_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a >= 0 AND b <= 100))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_compound_40e1387d_to.t
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
|   CONSTRAINT `t_OBCHECK_0000000000000007` CHECK (((`a` >= 0) and (`b` <= 100)))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_compound_40e1387d
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_compound_40e1387d DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_compound_40e1387d
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_compound_40e1387d`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_compound_40e1387d`.`t` ADD CHECK (((`a` >= 0) and (`b` <= 100)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_compound_40e1387d.t
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
|   CONSTRAINT `t_OBCHECK_0000000000000008` CHECK (((`a` >= 0) and (`b` <= 100)))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_string_literal_fc2e6593_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_string_literal_fc2e6593_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_string_literal_fc2e6593_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(20), CHECK (s = 'hello')) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_string_literal_fc2e6593_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_OBCHECK_0000000000000009` CHECK ((`s` = 'hello'))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_string_literal_fc2e6593
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_string_literal_fc2e6593 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_string_literal_fc2e6593
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_string_literal_fc2e6593`.`t` (
|   `id` int NOT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_string_literal_fc2e6593`.`t` ADD CHECK ((`s` = 'hello'))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_string_literal_fc2e6593.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_OBCHECK_0000000000000010` CHECK ((`s` = 'hello'))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_function_90a261c8_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_function_90a261c8_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_function_90a261c8_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(20), CHECK (CHAR_LENGTH(s) > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_function_90a261c8_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_OBCHECK_0000000000000011` CHECK ((CHAR_LENGTH(`s`) > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_function_90a261c8
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_function_90a261c8 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_function_90a261c8
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_function_90a261c8`.`t` (
|   `id` int NOT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_function_90a261c8`.`t` ADD CHECK ((char_length(`s`) > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_function_90a261c8.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_OBCHECK_0000000000000012` CHECK ((char_length(`s`) > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_column_level_check_50be9ec8_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_column_level_check_50be9ec8_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_column_level_check_50be9ec8_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_column_level_check_50be9ec8_to.t
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
|   CONSTRAINT `t_OBCHECK_0000000000000013` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_create_with_column_level_check_50be9ec8
< ok 0 0
> exec
| CREATE DATABASE pa_create_with_column_level_check_50be9ec8 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_with_column_level_check_50be9ec8
< ok 0 0
> exec
| CREATE TABLE `pa_create_with_column_level_check_50be9ec8`.`t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_with_column_level_check_50be9ec8`.`t` ADD CHECK ((`a` > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_with_column_level_check_50be9ec8.t
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
|   CONSTRAINT `t_OBCHECK_0000000000000014` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_named_check_b237fb0c_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_named_check_b237fb0c_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_named_check_b237fb0c_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_named_check_b237fb0c_to.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_named_check_b237fb0c
< ok 0 0
> exec
| CREATE DATABASE pa_add_named_check_b237fb0c DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_named_check_b237fb0c
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_named_check_b237fb0c.t
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
| ALTER TABLE `pa_add_named_check_b237fb0c`.`t` ADD CONSTRAINT `ck` CHECK ((`a` > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_named_check_b237fb0c.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_unnamed_check_c4055889_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_unnamed_check_c4055889_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_unnamed_check_c4055889_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_unnamed_check_c4055889_to.t
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
|   CONSTRAINT `t_OBCHECK_0000000000000015` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_unnamed_check_c4055889
< ok 0 0
> exec
| CREATE DATABASE pa_add_unnamed_check_c4055889 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_unnamed_check_c4055889
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_unnamed_check_c4055889.t
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
| ALTER TABLE `pa_add_unnamed_check_c4055889`.`t` ADD CHECK ((`a` > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_unnamed_check_c4055889.t
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
|   CONSTRAINT `t_OBCHECK_0000000000000016` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_not_enforced_b9ec14c6_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_not_enforced_b9ec14c6_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_not_enforced_b9ec14c6_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ne CHECK (a < b) NOT ENFORCED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_not_enforced_b9ec14c6_to.t
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
|   CONSTRAINT `ne` CHECK ((`a` < `b`))  NOT ENFORCED
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_not_enforced_b9ec14c6
< ok 0 0
> exec
| CREATE DATABASE pa_add_not_enforced_b9ec14c6 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_not_enforced_b9ec14c6
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_not_enforced_b9ec14c6.t
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
| ALTER TABLE `pa_add_not_enforced_b9ec14c6`.`t` ADD CONSTRAINT `ne` CHECK ((`a` < `b`)) /*!80016 NOT ENFORCED */
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_not_enforced_b9ec14c6.t
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
|   CONSTRAINT `ne` CHECK ((`a` < `b`))  NOT ENFORCED
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_second_check_41a0f209_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_second_check_41a0f209_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_second_check_41a0f209_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT c1 CHECK (a > 0), CONSTRAINT c2 CHECK (b < 100))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_second_check_41a0f209_to.t
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
|   CONSTRAINT `c1` CHECK ((`a` > 0)),
|   CONSTRAINT `c2` CHECK ((`b` < 100))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_second_check_41a0f209
< ok 0 0
> exec
| CREATE DATABASE pa_add_second_check_41a0f209 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_second_check_41a0f209
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT c1 CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_second_check_41a0f209.t
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
|   CONSTRAINT `c1` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_add_second_check_41a0f209`.`t` ADD CONSTRAINT `c2` CHECK ((`b` < 100))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_second_check_41a0f209.t
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
|   CONSTRAINT `c1` CHECK ((`a` > 0)),
|   CONSTRAINT `c2` CHECK ((`b` < 100))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_compound_6d4bc981_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_compound_6d4bc981_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_compound_6d4bc981_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a >= 0 AND b <= 100))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_compound_6d4bc981_to.t
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
|   CONSTRAINT `t_OBCHECK_0000000000000017` CHECK (((`a` >= 0) and (`b` <= 100)))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_compound_6d4bc981
< ok 0 0
> exec
| CREATE DATABASE pa_add_compound_6d4bc981 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_compound_6d4bc981
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_compound_6d4bc981.t
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
| ALTER TABLE `pa_add_compound_6d4bc981`.`t` ADD CHECK (((`a` >= 0) and (`b` <= 100)))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_compound_6d4bc981.t
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
|   CONSTRAINT `t_OBCHECK_0000000000000018` CHECK (((`a` >= 0) and (`b` <= 100)))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_drop_named_check_b5f1e49a_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_named_check_b5f1e49a_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_named_check_b5f1e49a_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_named_check_b5f1e49a_to.t
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
| DROP DATABASE IF EXISTS pa_drop_named_check_b5f1e49a
< ok 0 0
> exec
| CREATE DATABASE pa_drop_named_check_b5f1e49a DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_named_check_b5f1e49a
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_named_check_b5f1e49a.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_drop_named_check_b5f1e49a`.`t` DROP CHECK `ck`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_named_check_b5f1e49a.t
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
| DROP DATABASE IF EXISTS pa_drop_unnamed_check_5b81f3c3_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_unnamed_check_5b81f3c3_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_unnamed_check_5b81f3c3_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_unnamed_check_5b81f3c3_to.t
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
| DROP DATABASE IF EXISTS pa_drop_unnamed_check_5b81f3c3
< ok 0 0
> exec
| CREATE DATABASE pa_drop_unnamed_check_5b81f3c3 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_unnamed_check_5b81f3c3
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_unnamed_check_5b81f3c3.t
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
|   CONSTRAINT `t_OBCHECK_0000000000000019` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_drop_unnamed_check_5b81f3c3`.`t` DROP CHECK `t_OBCHECK_0000000000000019`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_unnamed_check_5b81f3c3.t
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
| DROP DATABASE IF EXISTS pa_drop_one_of_many_c3ed9fdc_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_one_of_many_c3ed9fdc_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_one_of_many_c3ed9fdc_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT c2 CHECK (b < 100))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_one_of_many_c3ed9fdc_to.t
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
|   CONSTRAINT `c2` CHECK ((`b` < 100))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_drop_one_of_many_c3ed9fdc
< ok 0 0
> exec
| CREATE DATABASE pa_drop_one_of_many_c3ed9fdc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_one_of_many_c3ed9fdc
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT c1 CHECK (a > 0), CONSTRAINT c2 CHECK (b < 100))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_one_of_many_c3ed9fdc.t
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
|   CONSTRAINT `c1` CHECK ((`a` > 0)),
|   CONSTRAINT `c2` CHECK ((`b` < 100))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_drop_one_of_many_c3ed9fdc`.`t` DROP CHECK `c1`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_one_of_many_c3ed9fdc.t
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
|   CONSTRAINT `c2` CHECK ((`b` < 100))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_expr_93868a64_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_expr_93868a64_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_expr_93868a64_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 10))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_expr_93868a64_to.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 10))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_expr_93868a64
< ok 0 0
> exec
| CREATE DATABASE pa_modify_expr_93868a64 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_expr_93868a64
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_expr_93868a64.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_expr_93868a64`.`t` DROP CHECK `ck`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_expr_93868a64`.`t` ADD CONSTRAINT `ck` CHECK ((`a` > 10))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_expr_93868a64.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 10))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_enforced_to_not_6cb9ee3d_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_enforced_to_not_6cb9ee3d_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_enforced_to_not_6cb9ee3d_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b) NOT ENFORCED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_enforced_to_not_6cb9ee3d_to.t
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
|   CONSTRAINT `ck` CHECK ((`a` < `b`))  NOT ENFORCED
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_enforced_to_not_6cb9ee3d
< ok 0 0
> exec
| CREATE DATABASE pa_modify_enforced_to_not_6cb9ee3d DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_enforced_to_not_6cb9ee3d
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_enforced_to_not_6cb9ee3d.t
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
|   CONSTRAINT `ck` CHECK ((`a` < `b`))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_enforced_to_not_6cb9ee3d`.`t` DROP CHECK `ck`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_enforced_to_not_6cb9ee3d`.`t` ADD CONSTRAINT `ck` CHECK ((`a` < `b`)) /*!80016 NOT ENFORCED */
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_enforced_to_not_6cb9ee3d.t
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
|   CONSTRAINT `ck` CHECK ((`a` < `b`))  NOT ENFORCED
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_not_to_enforced_53b9597f_to
< ok 0 0
> exec
| CREATE DATABASE pa_modify_not_to_enforced_53b9597f_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_not_to_enforced_53b9597f_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_not_to_enforced_53b9597f_to.t
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
|   CONSTRAINT `ck` CHECK ((`a` < `b`))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_modify_not_to_enforced_53b9597f
< ok 0 0
> exec
| CREATE DATABASE pa_modify_not_to_enforced_53b9597f DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_modify_not_to_enforced_53b9597f
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b) NOT ENFORCED)
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_not_to_enforced_53b9597f.t
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
|   CONSTRAINT `ck` CHECK ((`a` < `b`))  NOT ENFORCED
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_modify_not_to_enforced_53b9597f`.`t` DROP CHECK `ck`
< ok 0 0
> exec
| ALTER TABLE `pa_modify_not_to_enforced_53b9597f`.`t` ADD CONSTRAINT `ck` CHECK ((`a` < `b`))
< ok 0 0
> query
| SHOW CREATE TABLE pa_modify_not_to_enforced_53b9597f.t
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
|   CONSTRAINT `ck` CHECK ((`a` < `b`))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_column_and_check_on_it_294b18a3_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_and_check_on_it_294b18a3_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_and_check_on_it_294b18a3_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_and_check_on_it_294b18a3_to.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_add_column_and_check_on_it_294b18a3
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_and_check_on_it_294b18a3 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_and_check_on_it_294b18a3
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_and_check_on_it_294b18a3.t
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
| ALTER TABLE `pa_add_column_and_check_on_it_294b18a3`.`t` ADD COLUMN `a` int DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `pa_add_column_and_check_on_it_294b18a3`.`t` ADD CONSTRAINT `ck` CHECK ((`a` > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_and_check_on_it_294b18a3.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS pa_drop_check_and_its_column_9a933b0e_to
< ok 0 0
> exec
| CREATE DATABASE pa_drop_check_and_its_column_9a933b0e_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_check_and_its_column_9a933b0e_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_check_and_its_column_9a933b0e_to.t
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
| DROP DATABASE IF EXISTS pa_drop_check_and_its_column_9a933b0e
< ok 0 0
> exec
| CREATE DATABASE pa_drop_check_and_its_column_9a933b0e DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_drop_check_and_its_column_9a933b0e
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_check_and_its_column_9a933b0e.t
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
|   CONSTRAINT `ck` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `pa_drop_check_and_its_column_9a933b0e`.`t` DROP CHECK `ck`
< ok 0 0
> exec
| ALTER TABLE `pa_drop_check_and_its_column_9a933b0e`.`t` DROP COLUMN `a`
< ok 0 0
> query
| SHOW CREATE TABLE pa_drop_check_and_its_column_9a933b0e.t
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
