-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_CrossObjectOrderingApplyCorrectness default#1
> exec
| DROP DATABASE IF EXISTS ordsdl_to_create_function_and_view
< ok 0 0
> exec
| CREATE DATABASE ordsdl_to_create_function_and_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ordsdl_to_create_function_and_view
< ok 0 0
> exec
| CREATE TABLE users (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE FUNCTION ent_ucount() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) FROM users)
< ok 0 0
> exec
| CREATE VIEW ent_v_fn AS SELECT ent_ucount() AS n
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_to_create_function_and_view.users
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW ordsdl_to_create_function_and_view.ent_v_fn
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| ent_v_fn
: bytes
| CREATE VIEW `ent_v_fn` AS select `ordsdl_to_create_function_and_view`.`ent_ucount`() AS `n`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE FUNCTION ordsdl_to_create_function_and_view.ent_ucount
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| ent_ucount
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `ent_ucount`() RETURNS int(11)
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM users)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS ordsdl_from_create_function_and_view
< ok 0 0
> exec
| CREATE DATABASE ordsdl_from_create_function_and_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ordsdl_from_create_function_and_view
< ok 0 0
> exec
| CREATE TABLE users (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_from_create_function_and_view.users
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS ordsdl_t
< ok 0 0
> exec
| CREATE DATABASE ordsdl_t DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ordsdl_t
< ok 0 0
> exec
| CREATE TABLE users (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE FUNCTION `ordsdl_t`.`ent_ucount`() RETURNS int(11)
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM users)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `ordsdl_t`.`ent_v_fn` AS select ordsdl_t.ent_ucount() AS `n`
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_t.users
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW ordsdl_t.ent_v_fn
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| ent_v_fn
: bytes
| CREATE VIEW `ent_v_fn` AS select `ordsdl_t`.`ent_ucount`() AS `n`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE FUNCTION ordsdl_t.ent_ucount
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| ent_ucount
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `ent_ucount`() RETURNS int(11)
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM users)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS ordsdl_t
< ok 2 0
> exec
| DROP DATABASE IF EXISTS ordsdl_to_create_function_and_view
< ok 2 0
> exec
| DROP DATABASE IF EXISTS ordsdl_from_create_function_and_view
< ok 1 0
> exec
| DROP DATABASE IF EXISTS ordsdl_to_drop_function_and_view
< ok 0 0
> exec
| CREATE DATABASE ordsdl_to_drop_function_and_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ordsdl_to_drop_function_and_view
< ok 0 0
> exec
| CREATE TABLE users (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_to_drop_function_and_view.users
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS ordsdl_from_drop_function_and_view
< ok 0 0
> exec
| CREATE DATABASE ordsdl_from_drop_function_and_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ordsdl_from_drop_function_and_view
< ok 0 0
> exec
| CREATE TABLE users (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE FUNCTION ent_ucount() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) FROM users)
< ok 0 0
> exec
| CREATE VIEW ent_v_fn AS SELECT ent_ucount() AS n
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_from_drop_function_and_view.users
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW ordsdl_from_drop_function_and_view.ent_v_fn
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| ent_v_fn
: bytes
| CREATE VIEW `ent_v_fn` AS select `ordsdl_from_drop_function_and_view`.`ent_ucount`() AS `n`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE FUNCTION ordsdl_from_drop_function_and_view.ent_ucount
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| ent_ucount
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `ent_ucount`() RETURNS int(11)
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM users)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS ordsdl_t
< ok 0 0
> exec
| CREATE DATABASE ordsdl_t DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ordsdl_t
< ok 0 0
> exec
| CREATE TABLE users (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE FUNCTION ent_ucount() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) FROM users)
< ok 0 0
> exec
| CREATE VIEW ent_v_fn AS SELECT ent_ucount() AS n
< ok 0 0
> exec
| DROP VIEW IF EXISTS `ordsdl_t`.`ent_v_fn`
< ok 0 0
> exec
| DROP FUNCTION `ordsdl_t`.`ent_ucount`
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_t.users
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW ordsdl_t.ent_v_fn
< error
| 1146 42S02 Table 'ordsdl_t.ent_v_fn' doesn't exist
> query
| SHOW CREATE FUNCTION ordsdl_t.ent_ucount
< error
| 1305 42000 procedure/function does not exist
> exec
| DROP DATABASE IF EXISTS ordsdl_t
< ok 1 0
> exec
| DROP DATABASE IF EXISTS ordsdl_to_drop_function_and_view
< ok 1 0
> exec
| DROP DATABASE IF EXISTS ordsdl_from_drop_function_and_view
< ok 2 0
> exec
| DROP DATABASE IF EXISTS ordsdl_to_create_function_view_view_chain
< ok 0 0
> exec
| CREATE DATABASE ordsdl_to_create_function_view_view_chain DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ordsdl_to_create_function_view_view_chain
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE FUNCTION fn_base() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) FROM t)
< ok 0 0
> exec
| CREATE VIEW v1 AS SELECT fn_base() AS n
< ok 0 0
> exec
| CREATE VIEW v2 AS SELECT n FROM v1
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_to_create_function_view_view_chain.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW ordsdl_to_create_function_view_view_chain.v1
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v1
: bytes
| CREATE VIEW `v1` AS select `ordsdl_to_create_function_view_view_chain`.`fn_base`() AS `n`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW ordsdl_to_create_function_view_view_chain.v2
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v2
: bytes
| CREATE VIEW `v2` AS select `ordsdl_to_create_function_view_view_chain`.`v1`.`n` AS `n` from `ordsdl_to_create_function_view_view_chain`.`v1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE FUNCTION ordsdl_to_create_function_view_view_chain.fn_base
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_base
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `fn_base`() RETURNS int(11)
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM t)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS ordsdl_from_create_function_view_view_chain
< ok 0 0
> exec
| CREATE DATABASE ordsdl_from_create_function_view_view_chain DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ordsdl_from_create_function_view_view_chain
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_from_create_function_view_view_chain.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS ordsdl_t
< ok 0 0
> exec
| CREATE DATABASE ordsdl_t DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ordsdl_t
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE FUNCTION `ordsdl_t`.`fn_base`() RETURNS int(11)
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM t)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `ordsdl_t`.`v1` AS select ordsdl_t.fn_base() AS `n`
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `ordsdl_t`.`v2` AS select `ordsdl_t`.`v1`.`n` AS `n` from `ordsdl_t`.`v1`
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_t.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW ordsdl_t.v1
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v1
: bytes
| CREATE VIEW `v1` AS select `ordsdl_t`.`fn_base`() AS `n`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW ordsdl_t.v2
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v2
: bytes
| CREATE VIEW `v2` AS select `ordsdl_t`.`v1`.`n` AS `n` from `ordsdl_t`.`v1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE FUNCTION ordsdl_t.fn_base
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_base
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `fn_base`() RETURNS int(11)
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM t)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS ordsdl_t
< ok 3 0
> exec
| DROP DATABASE IF EXISTS ordsdl_to_create_function_view_view_chain
< ok 3 0
> exec
| DROP DATABASE IF EXISTS ordsdl_from_create_function_view_view_chain
< ok 1 0
> exec
| DROP DATABASE IF EXISTS ordsdl_to_create_function_selecting_view
< ok 0 0
> exec
| CREATE DATABASE ordsdl_to_create_function_selecting_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ordsdl_to_create_function_selecting_view
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v_src AS SELECT a FROM t
< ok 0 0
> exec
| CREATE FUNCTION fn_over_view() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) FROM v_src)
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_to_create_function_selecting_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW ordsdl_to_create_function_selecting_view.v_src
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_src
: bytes
| CREATE VIEW `v_src` AS select `ordsdl_to_create_function_selecting_view`.`t`.`a` AS `a` from `ordsdl_to_create_function_selecting_view`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE FUNCTION ordsdl_to_create_function_selecting_view.fn_over_view
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_over_view
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `fn_over_view`() RETURNS int(11)
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM v_src)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS ordsdl_from_create_function_selecting_view
< ok 0 0
> exec
| CREATE DATABASE ordsdl_from_create_function_selecting_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ordsdl_from_create_function_selecting_view
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_from_create_function_selecting_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS ordsdl_t
< ok 0 0
> exec
| CREATE DATABASE ordsdl_t DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ordsdl_t
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE FUNCTION `ordsdl_t`.`fn_over_view`() RETURNS int(11)
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM v_src)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `ordsdl_t`.`v_src` AS select `ordsdl_t`.`t`.`a` AS `a` from `ordsdl_t`.`t`
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_t.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW ordsdl_t.v_src
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_src
: bytes
| CREATE VIEW `v_src` AS select `ordsdl_t`.`t`.`a` AS `a` from `ordsdl_t`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE FUNCTION ordsdl_t.fn_over_view
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_over_view
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `fn_over_view`() RETURNS int(11)
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM v_src)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS ordsdl_t
< ok 2 0
> exec
| DROP DATABASE IF EXISTS ordsdl_to_create_function_selecting_view
< ok 2 0
> exec
| DROP DATABASE IF EXISTS ordsdl_from_create_function_selecting_view
< ok 1 0
> exec
| DROP DATABASE IF EXISTS ordsdl_to_modify_function_add_view
< ok 0 0
> exec
| CREATE DATABASE ordsdl_to_modify_function_add_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ordsdl_to_modify_function_add_view
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE FUNCTION fn_base() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) + 1 FROM t)
< ok 0 0
> exec
| CREATE VIEW v1 AS SELECT fn_base() AS n
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_to_modify_function_add_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW ordsdl_to_modify_function_add_view.v1
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v1
: bytes
| CREATE VIEW `v1` AS select `ordsdl_to_modify_function_add_view`.`fn_base`() AS `n`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE FUNCTION ordsdl_to_modify_function_add_view.fn_base
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_base
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `fn_base`() RETURNS int(11)
|     READS SQL DATA
| RETURN (SELECT COUNT(*) + 1 FROM t)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS ordsdl_from_modify_function_add_view
< ok 0 0
> exec
| CREATE DATABASE ordsdl_from_modify_function_add_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ordsdl_from_modify_function_add_view
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE FUNCTION fn_base() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) FROM t)
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_from_modify_function_add_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE FUNCTION ordsdl_from_modify_function_add_view.fn_base
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_base
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `fn_base`() RETURNS int(11)
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM t)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS ordsdl_t
< ok 0 0
> exec
| CREATE DATABASE ordsdl_t DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE ordsdl_t
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE FUNCTION fn_base() RETURNS INT READS SQL DATA RETURN (SELECT COUNT(*) FROM t)
< ok 0 0
> exec
| DROP FUNCTION `ordsdl_t`.`fn_base`
< ok 0 0
> exec
| CREATE FUNCTION `ordsdl_t`.`fn_base`() RETURNS int(11)
|     READS SQL DATA
| RETURN (SELECT COUNT(*) + 1 FROM t)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `ordsdl_t`.`v1` AS select ordsdl_t.fn_base() AS `n`
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_t.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW ordsdl_t.v1
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v1
: bytes
| CREATE VIEW `v1` AS select `ordsdl_t`.`fn_base`() AS `n`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE FUNCTION ordsdl_t.fn_base
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_base
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `fn_base`() RETURNS int(11)
|     READS SQL DATA
| RETURN (SELECT COUNT(*) + 1 FROM t)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS ordsdl_t
< ok 2 0
> exec
| DROP DATABASE IF EXISTS ordsdl_to_modify_function_add_view
< ok 2 0
> exec
| DROP DATABASE IF EXISTS ordsdl_from_modify_function_add_view
< ok 1 0
