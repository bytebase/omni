-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_CrossObjectOrderingApplyCorrectness default#1
> query
| SHOW CREATE FUNCTION ordsdl_to_create_function_and_view.ent_ucount
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| ent_ucount
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `ent_ucount`() RETURNS int
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM users)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE FUNCTION ordsdl_t.ent_ucount
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| ent_ucount
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `ent_ucount`() RETURNS int
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM users)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ordsdl_t
< ok 2 0
> exec
| DROP DATABASE IF EXISTS ordsdl_to_create_function_and_view
< ok 2 0
> exec
| DROP DATABASE IF EXISTS ordsdl_from_create_function_and_view
< ok 1 0
> query
| SHOW CREATE FUNCTION ordsdl_from_drop_function_and_view.ent_ucount
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| ent_ucount
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `ent_ucount`() RETURNS int
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM users)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ordsdl_t
< ok 1 0
> exec
| DROP DATABASE IF EXISTS ordsdl_to_drop_function_and_view
< ok 1 0
> exec
| DROP DATABASE IF EXISTS ordsdl_from_drop_function_and_view
< ok 2 0
> query
| SHOW CREATE FUNCTION ordsdl_to_create_function_view_view_chain.fn_base
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_base
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_base`() RETURNS int
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM t)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE FUNCTION ordsdl_t.fn_base
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_base
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_base`() RETURNS int
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM t)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ordsdl_t
< ok 3 0
> exec
| DROP DATABASE IF EXISTS ordsdl_to_create_function_view_view_chain
< ok 3 0
> exec
| DROP DATABASE IF EXISTS ordsdl_from_create_function_view_view_chain
< ok 1 0
> query
| SHOW CREATE FUNCTION ordsdl_to_create_function_selecting_view.fn_over_view
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_over_view
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_over_view`() RETURNS int
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM v_src)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE FUNCTION ordsdl_t.fn_over_view
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_over_view
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_over_view`() RETURNS int
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM v_src)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ordsdl_t
< ok 2 0
> exec
| DROP DATABASE IF EXISTS ordsdl_to_create_function_selecting_view
< ok 2 0
> exec
| DROP DATABASE IF EXISTS ordsdl_from_create_function_selecting_view
< ok 1 0
> query
| SHOW CREATE FUNCTION ordsdl_to_modify_function_add_view.fn_base
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_base
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_base`() RETURNS int
|     READS SQL DATA
| RETURN (SELECT COUNT(*) + 1 FROM t)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE FUNCTION ordsdl_from_modify_function_add_view.fn_base
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_base
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_base`() RETURNS int
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM t)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE FUNCTION ordsdl_t.fn_base
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_base
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_base`() RETURNS int
|     READS SQL DATA
| RETURN (SELECT COUNT(*) + 1 FROM t)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS ordsdl_t
< ok 2 0
> exec
| DROP DATABASE IF EXISTS ordsdl_to_modify_function_add_view
< ok 2 0
> exec
| DROP DATABASE IF EXISTS ordsdl_from_modify_function_add_view
< ok 1 0

== TestOracle_CrossObjectOrderingApplyCorrectness/8.0/create-function-and-view default#1
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
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW ordsdl_to_create_function_and_view.ent_v_fn
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| ent_v_fn
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `ent_v_fn` AS select `ent_ucount`() AS `n`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_CrossObjectOrderingApplyCorrectness/8.0/create-function-and-view default#2
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
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CrossObjectOrderingApplyCorrectness/8.0/create-function-and-view default#3
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
| CREATE FUNCTION `ordsdl_t`.`ent_ucount`() RETURNS int
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM users)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `ordsdl_t`.`ent_v_fn` AS select ent_ucount() AS `n`
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_t.users
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW ordsdl_t.ent_v_fn
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| ent_v_fn
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `ent_v_fn` AS select `ent_ucount`() AS `n`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_CrossObjectOrderingApplyCorrectness/8.0/create-function-selecting-view default#1
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW ordsdl_to_create_function_selecting_view.v_src
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_src
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_src` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_CrossObjectOrderingApplyCorrectness/8.0/create-function-selecting-view default#2
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CrossObjectOrderingApplyCorrectness/8.0/create-function-selecting-view default#3
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
| CREATE FUNCTION `ordsdl_t`.`fn_over_view`() RETURNS int
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM v_src)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `ordsdl_t`.`v_src` AS select `t`.`a` AS `a` from `t`
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_t.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW ordsdl_t.v_src
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_src
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_src` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_CrossObjectOrderingApplyCorrectness/8.0/create-function-view-view-chain default#1
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW ordsdl_to_create_function_view_view_chain.v1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v1` AS select `fn_base`() AS `n`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW ordsdl_to_create_function_view_view_chain.v2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v2` AS select `v1`.`n` AS `n` from `v1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_CrossObjectOrderingApplyCorrectness/8.0/create-function-view-view-chain default#2
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CrossObjectOrderingApplyCorrectness/8.0/create-function-view-view-chain default#3
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
| CREATE FUNCTION `ordsdl_t`.`fn_base`() RETURNS int
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM t)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `ordsdl_t`.`v1` AS select fn_base() AS `n`
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `ordsdl_t`.`v2` AS select `v1`.`n` AS `n` from `v1`
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_t.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW ordsdl_t.v1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v1` AS select `fn_base`() AS `n`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW ordsdl_t.v2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v2` AS select `v1`.`n` AS `n` from `v1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_CrossObjectOrderingApplyCorrectness/8.0/drop-function-and-view default#1
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
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CrossObjectOrderingApplyCorrectness/8.0/drop-function-and-view default#2
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
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW ordsdl_from_drop_function_and_view.ent_v_fn
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| ent_v_fn
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `ent_v_fn` AS select `ent_ucount`() AS `n`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_CrossObjectOrderingApplyCorrectness/8.0/drop-function-and-view default#3
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
< column "Create Table" "VARCHAR"
< row
: bytes
| users
: bytes
| CREATE TABLE `users` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW ordsdl_t.ent_v_fn
< error
| 1146 42S02 Table 'ordsdl_t.ent_v_fn' doesn't exist
> query
| SHOW CREATE FUNCTION ordsdl_t.ent_ucount
< error
| 1305 42000 FUNCTION ent_ucount does not exist

== TestOracle_CrossObjectOrderingApplyCorrectness/8.0/modify-function-add-view default#1
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW ordsdl_to_modify_function_add_view.v1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v1` AS select `fn_base`() AS `n`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_CrossObjectOrderingApplyCorrectness/8.0/modify-function-add-view default#2
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_CrossObjectOrderingApplyCorrectness/8.0/modify-function-add-view default#3
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
| CREATE FUNCTION `ordsdl_t`.`fn_base`() RETURNS int
|     READS SQL DATA
| RETURN (SELECT COUNT(*) + 1 FROM t)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `ordsdl_t`.`v1` AS select fn_base() AS `n`
< ok 0 0
> query
| SHOW CREATE TABLE ordsdl_t.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW ordsdl_t.v1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v1` AS select `fn_base`() AS `n`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
