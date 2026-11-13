-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_ViewIdempotence default#1
> exec
| DROP DATABASE IF EXISTS videm_simple
< ok 0 0
> exec
| CREATE DATABASE videm_simple DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_simple
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20), c INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a, b FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_simple.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_simple.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_simple`.`t`.`a` AS `a`,`videm_simple`.`t`.`b` AS `b` from `videm_simple`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_star_expanded
< ok 0 0
> exec
| CREATE DATABASE videm_star_expanded DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_star_expanded
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20), c INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT * FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_star_expanded.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_star_expanded.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_star_expanded`.`t`.`a` AS `a`,`videm_star_expanded`.`t`.`b` AS `b`,`videm_star_expanded`.`t`.`c` AS `c` from `videm_star_expanded`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_where_expr
< ok 0 0
> exec
| CREATE DATABASE videm_where_expr DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_where_expr
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t WHERE a > 0
< ok 0 0
> query
| SHOW CREATE TABLE videm_where_expr.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_where_expr.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_where_expr`.`t`.`a` AS `a` from `videm_where_expr`.`t` where (`videm_where_expr`.`t`.`a` > 0)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_expr_and_func
< ok 0 0
> exec
| CREATE DATABASE videm_expr_and_func DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_expr_and_func
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a + 1 AS s, UPPER(b) AS ub FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_expr_and_func.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_expr_and_func.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select (`videm_expr_and_func`.`t`.`a` + 1) AS `s`,UPPER(`videm_expr_and_func`.`t`.`b`) AS `ub` from `videm_expr_and_func`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_join
< ok 0 0
> exec
| CREATE DATABASE videm_join DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_join
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20), c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.a, u.y FROM t JOIN u ON t.c = u.x
< ok 0 0
> query
| SHOW CREATE TABLE videm_join.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_join.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int(11) DEFAULT NULL,
|   `y` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_join.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_join`.`t`.`a` AS `a`,`videm_join`.`u`.`y` AS `y` from (`videm_join`.`t` join `videm_join`.`u` on (`videm_join`.`t`.`c` = `videm_join`.`u`.`x`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_left_join
< ok 0 0
> exec
| CREATE DATABASE videm_left_join DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_left_join
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.a, u.y FROM t LEFT JOIN u ON t.c = u.x
< ok 0 0
> query
| SHOW CREATE TABLE videm_left_join.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_left_join.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int(11) DEFAULT NULL,
|   `y` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_left_join.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_left_join`.`t`.`a` AS `a`,`videm_left_join`.`u`.`y` AS `y` from (`videm_left_join`.`t` left join `videm_left_join`.`u` on (`videm_left_join`.`t`.`c` = `videm_left_join`.`u`.`x`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_join_3_table
< ok 0 0
> exec
| CREATE DATABASE videm_join_3_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_join_3_table
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.a, u.y, w.n FROM t JOIN u ON t.c = u.x JOIN w ON u.y = w.m
< ok 0 0
> query
| SHOW CREATE TABLE videm_join_3_table.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_join_3_table.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int(11) DEFAULT NULL,
|   `y` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_join_3_table.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int(11) DEFAULT NULL,
|   `n` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_join_3_table.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_join_3_table`.`t`.`a` AS `a`,`videm_join_3_table`.`u`.`y` AS `y`,`videm_join_3_table`.`w`.`n` AS `n` from ((`videm_join_3_table`.`t` join `videm_join_3_table`.`u` on (`videm_join_3_table`.`t`.`c` = `videm_join_3_table`.`u`.`x`)) join `videm_join_3_table`.`w` on (`videm_join_3_table`.`u`.`y` = `videm_join_3_table`.`w`.`m`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_join_4_table
< ok 0 0
> exec
| CREATE DATABASE videm_join_4_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_join_4_table
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> exec
| CREATE TABLE z (p INT, q INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.a FROM t JOIN u ON t.c = u.x JOIN w ON u.y = w.m JOIN z ON w.n = z.p
< ok 0 0
> query
| SHOW CREATE TABLE videm_join_4_table.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_join_4_table.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int(11) DEFAULT NULL,
|   `y` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_join_4_table.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int(11) DEFAULT NULL,
|   `n` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_join_4_table.z
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| z
: bytes
| CREATE TABLE `z` (
|   `p` int(11) DEFAULT NULL,
|   `q` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_join_4_table.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_join_4_table`.`t`.`a` AS `a` from (((`videm_join_4_table`.`t` join `videm_join_4_table`.`u` on (`videm_join_4_table`.`t`.`c` = `videm_join_4_table`.`u`.`x`)) join `videm_join_4_table`.`w` on (`videm_join_4_table`.`u`.`y` = `videm_join_4_table`.`w`.`m`)) join `videm_join_4_table`.`z` on (`videm_join_4_table`.`w`.`n` = `videm_join_4_table`.`z`.`p`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_join_left_inner_mix
< ok 0 0
> exec
| CREATE DATABASE videm_join_left_inner_mix DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_join_left_inner_mix
< ok 0 0
> exec
| CREATE TABLE category (category_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_category (film_id INT, category_id INT)
< ok 0 0
> exec
| CREATE TABLE film (film_id INT, title VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_actor (actor_id INT, film_id INT)
< ok 0 0
> exec
| CREATE TABLE actor (actor_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT film.film_id AS FID, category.name AS category, actor.name AS actor FROM category LEFT JOIN film_category ON category.category_id = film_category.category_id LEFT JOIN film ON film_category.film_id = film.film_id JOIN film_actor ON film.film_id = film_actor.film_id JOIN actor ON film_actor.actor_id = actor.actor_id
< ok 0 0
> query
| SHOW CREATE TABLE videm_join_left_inner_mix.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `category_id` int(11) DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_join_left_inner_mix.film_category
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| film_category
: bytes
| CREATE TABLE `film_category` (
|   `film_id` int(11) DEFAULT NULL,
|   `category_id` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_join_left_inner_mix.film
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| film
: bytes
| CREATE TABLE `film` (
|   `film_id` int(11) DEFAULT NULL,
|   `title` varchar(50) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_join_left_inner_mix.film_actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| film_actor
: bytes
| CREATE TABLE `film_actor` (
|   `actor_id` int(11) DEFAULT NULL,
|   `film_id` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_join_left_inner_mix.actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| actor
: bytes
| CREATE TABLE `actor` (
|   `actor_id` int(11) DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_join_left_inner_mix.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_join_left_inner_mix`.`film`.`film_id` AS `FID`,`videm_join_left_inner_mix`.`category`.`name` AS `category`,`videm_join_left_inner_mix`.`actor`.`name` AS `actor` from ((((`videm_join_left_inner_mix`.`category` left join `videm_join_left_inner_mix`.`film_category` on (`videm_join_left_inner_mix`.`category`.`category_id` = `videm_join_left_inner_mix`.`film_category`.`category_id`)) left join `videm_join_left_inner_mix`.`film` on (`videm_join_left_inner_mix`.`film_category`.`film_id` = `videm_join_left_inner_mix`.`film`.`film_id`)) join `videm_join_left_inner_mix`.`film_actor` on (`videm_join_left_inner_mix`.`film`.`film_id` = `videm_join_left_inner_mix`.`film_actor`.`film_id`)) join `videm_join_left_inner_mix`.`actor` on (`videm_join_left_inner_mix`.`film_actor`.`actor_id` = `videm_join_left_inner_mix`.`actor`.`actor_id`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_join_group_by
< ok 0 0
> exec
| CREATE DATABASE videm_join_group_by DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_join_group_by
< ok 0 0
> exec
| CREATE TABLE category (category_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_category (film_id INT, category_id INT)
< ok 0 0
> exec
| CREATE TABLE film (film_id INT, title VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_actor (actor_id INT, film_id INT)
< ok 0 0
> exec
| CREATE TABLE actor (actor_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT film.film_id AS FID, COUNT(*) AS n FROM category LEFT JOIN film_category ON category.category_id = film_category.category_id LEFT JOIN film ON film_category.film_id = film.film_id JOIN film_actor ON film.film_id = film_actor.film_id JOIN actor ON film_actor.actor_id = actor.actor_id GROUP BY film.film_id
< ok 0 0
> query
| SHOW CREATE TABLE videm_join_group_by.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `category_id` int(11) DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_join_group_by.film_category
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| film_category
: bytes
| CREATE TABLE `film_category` (
|   `film_id` int(11) DEFAULT NULL,
|   `category_id` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_join_group_by.film
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| film
: bytes
| CREATE TABLE `film` (
|   `film_id` int(11) DEFAULT NULL,
|   `title` varchar(50) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_join_group_by.film_actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| film_actor
: bytes
| CREATE TABLE `film_actor` (
|   `actor_id` int(11) DEFAULT NULL,
|   `film_id` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_join_group_by.actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| actor
: bytes
| CREATE TABLE `actor` (
|   `actor_id` int(11) DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_join_group_by.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_join_group_by`.`film`.`film_id` AS `FID`,count(*) AS `n` from ((((`videm_join_group_by`.`category` left join `videm_join_group_by`.`film_category` on (`videm_join_group_by`.`category`.`category_id` = `videm_join_group_by`.`film_category`.`category_id`)) left join `videm_join_group_by`.`film` on (`videm_join_group_by`.`film_category`.`film_id` = `videm_join_group_by`.`film`.`film_id`)) join `videm_join_group_by`.`film_actor` on (`videm_join_group_by`.`film`.`film_id` = `videm_join_group_by`.`film_actor`.`film_id`)) join `videm_join_group_by`.`actor` on (`videm_join_group_by`.`film_actor`.`actor_id` = `videm_join_group_by`.`actor`.`actor_id`)) group by `videm_join_group_by`.`film`.`film_id`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_view_on_multi_join_view
< ok 0 0
> exec
| CREATE DATABASE videm_view_on_multi_join_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_view_on_multi_join_view
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> exec
| CREATE VIEW base AS SELECT t.a AS a, u.y AS y, w.n AS n FROM t JOIN u ON t.c = u.x JOIN w ON u.y = w.m
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a, n FROM base WHERE a > 0
< ok 0 0
> query
| SHOW CREATE TABLE videm_view_on_multi_join_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_view_on_multi_join_view.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int(11) DEFAULT NULL,
|   `y` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_view_on_multi_join_view.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int(11) DEFAULT NULL,
|   `n` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_view_on_multi_join_view.base
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| base
: bytes
| CREATE VIEW `base` AS select `videm_view_on_multi_join_view`.`t`.`a` AS `a`,`videm_view_on_multi_join_view`.`u`.`y` AS `y`,`videm_view_on_multi_join_view`.`w`.`n` AS `n` from ((`videm_view_on_multi_join_view`.`t` join `videm_view_on_multi_join_view`.`u` on (`videm_view_on_multi_join_view`.`t`.`c` = `videm_view_on_multi_join_view`.`u`.`x`)) join `videm_view_on_multi_join_view`.`w` on (`videm_view_on_multi_join_view`.`u`.`y` = `videm_view_on_multi_join_view`.`w`.`m`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW videm_view_on_multi_join_view.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_view_on_multi_join_view`.`base`.`a` AS `a`,`videm_view_on_multi_join_view`.`base`.`n` AS `n` from `videm_view_on_multi_join_view`.`base` where (`videm_view_on_multi_join_view`.`base`.`a` > 0)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_join_aliased_compound_on
< ok 0 0
> exec
| CREATE DATABASE videm_join_aliased_compound_on DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_join_aliased_compound_on
< ok 0 0
> exec
| CREATE TABLE dept_emp (emp_no INT NOT NULL, dept_no CHAR(4) NOT NULL, from_date DATE NOT NULL, to_date DATE NOT NULL, PRIMARY KEY (emp_no, dept_no))
< ok 0 0
> exec
| CREATE TABLE departments (dept_no CHAR(4) NOT NULL, dept_name VARCHAR(40) NOT NULL, PRIMARY KEY (dept_no))
< ok 0 0
> exec
| CREATE VIEW dept_emp_latest_date AS SELECT emp_no, MAX(from_date) AS from_date, MAX(to_date) AS to_date FROM dept_emp GROUP BY emp_no
< ok 0 0
> exec
| CREATE VIEW current_dept_emp AS SELECT l.emp_no, d.dept_no, l.from_date, l.to_date FROM dept_emp d JOIN dept_emp_latest_date l ON d.emp_no = l.emp_no AND d.from_date = l.from_date JOIN departments dp ON d.dept_no = dp.dept_no
< ok 0 0
> query
| SHOW CREATE TABLE videm_join_aliased_compound_on.dept_emp
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| dept_emp
: bytes
| CREATE TABLE `dept_emp` (
|   `emp_no` int(11) NOT NULL,
|   `dept_no` char(4) NOT NULL,
|   `from_date` date NOT NULL,
|   `to_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`, `dept_no`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_join_aliased_compound_on.departments
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| departments
: bytes
| CREATE TABLE `departments` (
|   `dept_no` char(4) NOT NULL,
|   `dept_name` varchar(40) NOT NULL,
|   PRIMARY KEY (`dept_no`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_join_aliased_compound_on.dept_emp_latest_date
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| dept_emp_latest_date
: bytes
| CREATE VIEW `dept_emp_latest_date` AS select `videm_join_aliased_compound_on`.`dept_emp`.`emp_no` AS `emp_no`,max(`videm_join_aliased_compound_on`.`dept_emp`.`from_date`) AS `from_date`,max(`videm_join_aliased_compound_on`.`dept_emp`.`to_date`) AS `to_date` from `videm_join_aliased_compound_on`.`dept_emp` group by `videm_join_aliased_compound_on`.`dept_emp`.`emp_no`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW videm_join_aliased_compound_on.current_dept_emp
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| current_dept_emp
: bytes
| CREATE VIEW `current_dept_emp` AS select `videm_join_aliased_compound_on`.`l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`videm_join_aliased_compound_on`.`l`.`from_date` AS `from_date`,`videm_join_aliased_compound_on`.`l`.`to_date` AS `to_date` from ((`videm_join_aliased_compound_on`.`dept_emp` `d` join `videm_join_aliased_compound_on`.`dept_emp_latest_date` `l` on (`d`.`emp_no` = `videm_join_aliased_compound_on`.`l`.`emp_no`) and (`d`.`from_date` = `videm_join_aliased_compound_on`.`l`.`from_date`)) join `videm_join_aliased_compound_on`.`departments` `dp` on (`d`.`dept_no` = `dp`.`dept_no`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_join_aliased_left_inner_compound_on
< ok 0 0
> exec
| CREATE DATABASE videm_join_aliased_left_inner_compound_on DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_join_aliased_left_inner_compound_on
< ok 0 0
> exec
| CREATE TABLE a (id INT, k INT, v INT)
< ok 0 0
> exec
| CREATE TABLE b (id INT, k INT, w INT)
< ok 0 0
> exec
| CREATE TABLE c (id INT, k INT, z INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT x.v, y.w, q.z FROM a x LEFT JOIN b y ON x.id = y.id AND x.k = y.k INNER JOIN c q ON y.id = q.id AND y.k = q.k
< ok 0 0
> query
| SHOW CREATE TABLE videm_join_aliased_left_inner_compound_on.a
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `id` int(11) DEFAULT NULL,
|   `k` int(11) DEFAULT NULL,
|   `v` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_join_aliased_left_inner_compound_on.b
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| b
: bytes
| CREATE TABLE `b` (
|   `id` int(11) DEFAULT NULL,
|   `k` int(11) DEFAULT NULL,
|   `w` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_join_aliased_left_inner_compound_on.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) DEFAULT NULL,
|   `k` int(11) DEFAULT NULL,
|   `z` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_join_aliased_left_inner_compound_on.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `x`.`v` AS `v`,`y`.`w` AS `w`,`q`.`z` AS `z` from ((`videm_join_aliased_left_inner_compound_on`.`a` `x` left join `videm_join_aliased_left_inner_compound_on`.`b` `y` on (`x`.`id` = `y`.`id`) and (`x`.`k` = `y`.`k`)) join `videm_join_aliased_left_inner_compound_on`.`c` `q` on (`y`.`id` = `q`.`id`) and (`y`.`k` = `q`.`k`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_derived_table_from
< ok 0 0
> exec
| CREATE DATABASE videm_derived_table_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_derived_table_from
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT d.a FROM (SELECT a, b FROM t WHERE b > 0) d
< ok 0 0
> query
| SHOW CREATE TABLE videm_derived_table_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_derived_table_from.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `d`.`a` AS `a` from (select `videm_derived_table_from`.`t`.`a` AS `a`,`videm_derived_table_from`.`t`.`b` AS `b` from `videm_derived_table_from`.`t` where (`videm_derived_table_from`.`t`.`b` > 0)) `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_derived_table_join
< ok 0 0
> exec
| CREATE DATABASE videm_derived_table_join DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_derived_table_join
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.a, d.y FROM t JOIN (SELECT x, y FROM u) d ON t.c = d.x
< ok 0 0
> query
| SHOW CREATE TABLE videm_derived_table_join.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE videm_derived_table_join.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int(11) DEFAULT NULL,
|   `y` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_derived_table_join.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_derived_table_join`.`t`.`a` AS `a`,`d`.`y` AS `y` from (`videm_derived_table_join`.`t` join (select `videm_derived_table_join`.`u`.`x` AS `x`,`videm_derived_table_join`.`u`.`y` AS `y` from `videm_derived_table_join`.`u`) `d` on (`videm_derived_table_join`.`t`.`c` = `d`.`x`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_algorithm_merge
< ok 0 0
> exec
| CREATE DATABASE videm_algorithm_merge DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_algorithm_merge
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE ALGORITHM=MERGE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_algorithm_merge.t
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
| SHOW CREATE VIEW videm_algorithm_merge.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_algorithm_merge`.`t`.`a` AS `a` from `videm_algorithm_merge`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_algorithm_temptable
< ok 0 0
> exec
| CREATE DATABASE videm_algorithm_temptable DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_algorithm_temptable
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE ALGORITHM=TEMPTABLE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_algorithm_temptable.t
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
| SHOW CREATE VIEW videm_algorithm_temptable.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_algorithm_temptable`.`t`.`a` AS `a` from `videm_algorithm_temptable`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_sql_security_invoker
< ok 0 0
> exec
| CREATE DATABASE videm_sql_security_invoker DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_sql_security_invoker
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE SQL SECURITY INVOKER VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_sql_security_invoker.t
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
| SHOW CREATE VIEW videm_sql_security_invoker.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_sql_security_invoker`.`t`.`a` AS `a` from `videm_sql_security_invoker`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_check_cascaded
< ok 0 0
> exec
| CREATE DATABASE videm_check_cascaded DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_check_cascaded
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t WHERE a > 0 WITH CASCADED CHECK OPTION
< ok 0 0
> query
| SHOW CREATE TABLE videm_check_cascaded.t
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
| SHOW CREATE VIEW videm_check_cascaded.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_check_cascaded`.`t`.`a` AS `a` from `videm_check_cascaded`.`t` where (`videm_check_cascaded`.`t`.`a` > 0) WITH CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_check_local
< ok 0 0
> exec
| CREATE DATABASE videm_check_local DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_check_local
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t WHERE a > 0 WITH LOCAL CHECK OPTION
< ok 0 0
> query
| SHOW CREATE TABLE videm_check_local.t
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
| SHOW CREATE VIEW videm_check_local.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_check_local`.`t`.`a` AS `a` from `videm_check_local`.`t` where (`videm_check_local`.`t`.`a` > 0) WITH LOCAL CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_view_on_view
< ok 0 0
> exec
| CREATE DATABASE videm_view_on_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_view_on_view
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW base AS SELECT a, b FROM t
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM base
< ok 0 0
> query
| SHOW CREATE TABLE videm_view_on_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_view_on_view.base
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| base
: bytes
| CREATE VIEW `base` AS select `videm_view_on_view`.`t`.`a` AS `a`,`videm_view_on_view`.`t`.`b` AS `b` from `videm_view_on_view`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW videm_view_on_view.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_view_on_view`.`base`.`a` AS `a` from `videm_view_on_view`.`base`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_cte
< ok 0 0
> exec
| CREATE DATABASE videm_cte DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_cte
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS WITH cte AS (SELECT a, b FROM t) SELECT a FROM cte
< ok 0 0
> query
| SHOW CREATE TABLE videm_cte.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_cte.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS WITH RECURSIVE `cte` as (select `videm_cte`.`t`.`a` AS `a`,`videm_cte`.`t`.`b` AS `b` from `videm_cte`.`t`) select `cte`.`a` AS `a` from `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_union
< ok 0 0
> exec
| CREATE DATABASE videm_union DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_union
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE TABLE u (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t UNION SELECT a FROM u
< ok 0 0
> query
| SHOW CREATE TABLE videm_union.t
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
| SHOW CREATE TABLE videm_union.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_union.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS (select `videm_union`.`t`.`a` AS `a` from `videm_union`.`t`) union (select `videm_union`.`u`.`a` AS `a` from `videm_union`.`u`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_group_by_agg
< ok 0 0
> exec
| CREATE DATABASE videm_group_by_agg DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_group_by_agg
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a, COUNT(*) AS n FROM t GROUP BY a
< ok 0 0
> query
| SHOW CREATE TABLE videm_group_by_agg.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_group_by_agg.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_group_by_agg`.`t`.`a` AS `a`,count(*) AS `n` from `videm_group_by_agg`.`t` group by `videm_group_by_agg`.`t`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_explicit_columns
< ok 0 0
> exec
| CREATE DATABASE videm_explicit_columns DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_explicit_columns
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v (p, q) AS SELECT a, b FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_explicit_columns.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_explicit_columns.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_explicit_columns`.`t`.`a` AS `p`,`videm_explicit_columns`.`t`.`b` AS `q` from `videm_explicit_columns`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_view_on_view_explicit_cols
< ok 0 0
> exec
| CREATE DATABASE videm_view_on_view_explicit_cols DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_view_on_view_explicit_cols
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW base (x, y) AS SELECT a, b FROM t
< ok 0 0
> exec
| CREATE VIEW v AS SELECT x FROM base
< ok 0 0
> query
| SHOW CREATE TABLE videm_view_on_view_explicit_cols.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_view_on_view_explicit_cols.base
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| base
: bytes
| CREATE VIEW `base` AS select `videm_view_on_view_explicit_cols`.`t`.`a` AS `x`,`videm_view_on_view_explicit_cols`.`t`.`b` AS `y` from `videm_view_on_view_explicit_cols`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW videm_view_on_view_explicit_cols.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_view_on_view_explicit_cols`.`base`.`x` AS `x` from `videm_view_on_view_explicit_cols`.`base`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_compare_if
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_compare_if DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_compare_if
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT if(((SELECT count(0) FROM t) = 0),'NO','YES') AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_compare_if.t
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
| SHOW CREATE VIEW videm_paren_subquery_compare_if.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select (case when ((select count(0) from `videm_paren_subquery_compare_if`.`t`) = 0) then 'NO' else 'YES' end) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_compare_if_depth2
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_compare_if_depth2 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_compare_if_depth2
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT if((((SELECT count(0) FROM t)) = 0),'NO','YES') AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_compare_if_depth2.t
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
| SHOW CREATE VIEW videm_paren_subquery_compare_if_depth2.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select (case when ((select count(0) from `videm_paren_subquery_compare_if_depth2`.`t`) = 0) then 'NO' else 'YES' end) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_pure_depth2
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_pure_depth2 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_pure_depth2
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT ((SELECT count(0) FROM t)) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_pure_depth2.t
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
| SHOW CREATE VIEW videm_paren_subquery_pure_depth2.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select (select count(0) from `videm_paren_subquery_pure_depth2`.`t`) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_in_depth2
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_in_depth2 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_in_depth2
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT (1 IN ((SELECT max(a) FROM t))) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_in_depth2.t
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
| SHOW CREATE VIEW videm_paren_subquery_in_depth2.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select (1 = any (select max(`videm_paren_subquery_in_depth2`.`t`.`a`) from `videm_paren_subquery_in_depth2`.`t`)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_in_subquery
< ok 0 0
> exec
| CREATE DATABASE videm_in_subquery DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_in_subquery
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT (1 IN (SELECT max(a) FROM t)) AS x, (1 NOT IN (SELECT max(a) FROM t)) AS y
< ok 0 0
> query
| SHOW CREATE TABLE videm_in_subquery.t
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
| SHOW CREATE VIEW videm_in_subquery.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select (1 = any (select max(`videm_in_subquery`.`t`.`a`) from `videm_in_subquery`.`t`)) AS `x`,(1 <> all (select max(`videm_in_subquery`.`t`.`a`) from `videm_in_subquery`.`t`)) AS `y`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_in_subquery_where
< ok 0 0
> exec
| CREATE DATABASE videm_in_subquery_where DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_in_subquery_where
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE TABLE u (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t WHERE a IN (SELECT a FROM u) AND a NOT IN (SELECT a + 1 FROM u)
< ok 0 0
> query
| SHOW CREATE TABLE videm_in_subquery_where.t
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
| SHOW CREATE TABLE videm_in_subquery_where.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_in_subquery_where.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_in_subquery_where`.`t`.`a` AS `a` from `videm_in_subquery_where`.`t` where (`videm_in_subquery_where`.`t`.`a` = any (select `videm_in_subquery_where`.`u`.`a` from `videm_in_subquery_where`.`u`)) and (`videm_in_subquery_where`.`t`.`a` <> all (select (`videm_in_subquery_where`.`u`.`a` + 1) from `videm_in_subquery_where`.`u`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_union_chain
< ok 0 0
> exec
| CREATE DATABASE videm_union_chain DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_union_chain
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE TABLE u (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t UNION SELECT a FROM u UNION ALL SELECT 1
< ok 0 0
> query
| SHOW CREATE TABLE videm_union_chain.t
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
| SHOW CREATE TABLE videm_union_chain.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_union_chain.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS ((select `videm_union_chain`.`t`.`a` AS `a` from `videm_union_chain`.`t`) union (select `videm_union_chain`.`u`.`a` AS `a` from `videm_union_chain`.`u`)) union  all (select 1 AS `1`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_explicit_columns_union
< ok 0 0
> exec
| CREATE DATABASE videm_explicit_columns_union DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_explicit_columns_union
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE TABLE u (a INT)
< ok 0 0
> exec
| CREATE VIEW v (p) AS SELECT a FROM t UNION SELECT a FROM u
< ok 0 0
> query
| SHOW CREATE TABLE videm_explicit_columns_union.t
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
| SHOW CREATE TABLE videm_explicit_columns_union.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_explicit_columns_union.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS (select `videm_explicit_columns_union`.`t`.`a` AS `p` from `videm_explicit_columns_union`.`t`) union (select `videm_explicit_columns_union`.`u`.`a` AS `p` from `videm_explicit_columns_union`.`u`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_explicit_columns_expr
< ok 0 0
> exec
| CREATE DATABASE videm_explicit_columns_expr DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_explicit_columns_expr
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v (p, q) AS SELECT a, a + 1 FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_explicit_columns_expr.t
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
| SHOW CREATE VIEW videm_explicit_columns_expr.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_explicit_columns_expr`.`t`.`a` AS `p`,(`videm_explicit_columns_expr`.`t`.`a` + 1) AS `q` from `videm_explicit_columns_expr`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_order_by_alias
< ok 0 0
> exec
| CREATE DATABASE videm_order_by_alias DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_order_by_alias
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a AS x, b FROM t ORDER BY x DESC, b
< ok 0 0
> query
| SHOW CREATE TABLE videm_order_by_alias.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_order_by_alias.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_order_by_alias`.`t`.`a` AS `x`,`videm_order_by_alias`.`t`.`b` AS `b` from `videm_order_by_alias`.`t` order by 1 desc, 2
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_order_by_position
< ok 0 0
> exec
| CREATE DATABASE videm_order_by_position DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_order_by_position
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a, b FROM t ORDER BY 2, 1
< ok 0 0
> query
| SHOW CREATE TABLE videm_order_by_position.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_order_by_position.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_order_by_position`.`t`.`a` AS `a`,`videm_order_by_position`.`t`.`b` AS `b` from `videm_order_by_position`.`t` order by 2, 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_order_by_other_column
< ok 0 0
> exec
| CREATE DATABASE videm_order_by_other_column DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_order_by_other_column
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t ORDER BY b, a + 1
< ok 0 0
> query
| SHOW CREATE TABLE videm_order_by_other_column.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_order_by_other_column.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_order_by_other_column`.`t`.`a` AS `a` from `videm_order_by_other_column`.`t` order by `videm_order_by_other_column`.`t`.`b`, (`videm_order_by_other_column`.`t`.`a` + 1)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_order_by_union
< ok 0 0
> exec
| CREATE DATABASE videm_order_by_union DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_order_by_union
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE TABLE u (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t UNION SELECT a FROM u ORDER BY a LIMIT 2
< ok 0 0
> query
| SHOW CREATE TABLE videm_order_by_union.t
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
| SHOW CREATE TABLE videm_order_by_union.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_order_by_union.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS (select `videm_order_by_union`.`t`.`a` AS `a` from `videm_order_by_union`.`t`) union (select `videm_order_by_union`.`u`.`a` AS `a` from `videm_order_by_union`.`u`) order by 1 limit 2
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_explicit_columns_order_by_alias
< ok 0 0
> exec
| CREATE DATABASE videm_explicit_columns_order_by_alias DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_explicit_columns_order_by_alias
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v (p) AS SELECT a AS x FROM t ORDER BY x
< ok 0 0
> query
| SHOW CREATE TABLE videm_explicit_columns_order_by_alias.t
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
| SHOW CREATE VIEW videm_explicit_columns_order_by_alias.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_explicit_columns_order_by_alias`.`t`.`a` AS `p` from `videm_explicit_columns_order_by_alias`.`t` order by 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_explicit_columns_order_by_union
< ok 0 0
> exec
| CREATE DATABASE videm_explicit_columns_order_by_union DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_explicit_columns_order_by_union
< ok 0 0
> exec
| CREATE VIEW v (p) AS SELECT 1 AS a UNION SELECT 2 ORDER BY 1
< ok 0 0
> query
| SHOW CREATE VIEW videm_explicit_columns_order_by_union.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS (select 1 AS `p`) union (select 2 AS `p`) order by 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_explicit_columns_order_by_shadow
< ok 0 0
> exec
| CREATE DATABASE videm_explicit_columns_order_by_shadow DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_explicit_columns_order_by_shadow
< ok 0 0
> exec
| CREATE TABLE t (a INT, p INT)
< ok 0 0
> exec
| CREATE VIEW v (p) AS SELECT a AS x FROM t ORDER BY x + 1
< ok 0 0
> query
| SHOW CREATE TABLE videm_explicit_columns_order_by_shadow.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `p` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_explicit_columns_order_by_shadow.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_explicit_columns_order_by_shadow`.`t`.`a` AS `p` from `videm_explicit_columns_order_by_shadow`.`t` order by (`videm_explicit_columns_order_by_shadow`.`t`.`a` + 1)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_explicit_columns_order_by_convert
< ok 0 0
> exec
| CREATE DATABASE videm_explicit_columns_order_by_convert DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_explicit_columns_order_by_convert
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v (b) AS SELECT a AS x FROM t ORDER BY CONVERT(x, SIGNED)
< ok 0 0
> query
| SHOW CREATE TABLE videm_explicit_columns_order_by_convert.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_explicit_columns_order_by_convert.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_explicit_columns_order_by_convert`.`t`.`a` AS `b` from `videm_explicit_columns_order_by_convert`.`t` order by cast(`videm_explicit_columns_order_by_convert`.`t`.`a` as signed)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_explicit_columns_group_by
< ok 0 0
> exec
| CREATE DATABASE videm_explicit_columns_group_by DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_explicit_columns_group_by
< ok 0 0
> exec
| CREATE TABLE t (a INT, p INT)
< ok 0 0
> exec
| CREATE VIEW v (p, n) AS SELECT a AS x, count(*) AS c FROM t GROUP BY x HAVING c > 0 ORDER BY c DESC
< ok 0 0
> query
| SHOW CREATE TABLE videm_explicit_columns_group_by.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `p` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_explicit_columns_group_by.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_explicit_columns_group_by`.`t`.`a` AS `p`,count(*) AS `n` from `videm_explicit_columns_group_by`.`t` group by `videm_explicit_columns_group_by`.`t`.`a` having (count(*) > 0) order by 2 desc
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_paren_exists_depth2
< ok 0 0
> exec
| CREATE DATABASE videm_paren_exists_depth2 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_exists_depth2
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT (EXISTS((SELECT 1 FROM t))) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_exists_depth2.t
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
| SHOW CREATE VIEW videm_paren_exists_depth2.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select exists((select 1 from `videm_paren_exists_depth2`.`t`)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_setop_compare
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_setop_compare DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_setop_compare
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT ((SELECT 1 UNION SELECT 1) = 1) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_setop_compare.t
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
| SHOW CREATE VIEW videm_paren_subquery_setop_compare.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select (((select 1 AS `1`) union (select 1 AS `1`)) = 1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_arith
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_arith DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_arith
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT ((SELECT max(a) FROM t) + 1) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_arith.t
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
| SHOW CREATE VIEW videm_paren_subquery_arith.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select ((select max(`videm_paren_subquery_arith`.`t`.`a`) from `videm_paren_subquery_arith`.`t`) + 1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_where
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_where DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_where
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t WHERE ((SELECT count(0) FROM t) > 0)
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_where.t
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
| SHOW CREATE VIEW videm_paren_subquery_where.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_paren_subquery_where`.`t`.`a` AS `a` from `videm_paren_subquery_where`.`t` where ((select count(0) from `videm_paren_subquery_where`.`t`) > 0)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_both_operands
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_both_operands DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_both_operands
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT ((SELECT min(a) FROM t) = (SELECT max(a) FROM t)) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_both_operands.t
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
| SHOW CREATE VIEW videm_paren_subquery_both_operands.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select ((select min(`videm_paren_subquery_both_operands`.`t`.`a`) from `videm_paren_subquery_both_operands`.`t`) = (select max(`videm_paren_subquery_both_operands`.`t`.`a`) from `videm_paren_subquery_both_operands`.`t`)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_in_list
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_in_list DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_in_list
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT (1 IN ((SELECT max(a) FROM t) + 1, 2)) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_in_list.t
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
| SHOW CREATE VIEW videm_paren_subquery_in_list.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select (1 in (((select max(`videm_paren_subquery_in_list`.`t`.`a`) from `videm_paren_subquery_in_list`.`t`) + 1),2)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_paren_exists_compare
< ok 0 0
> exec
| CREATE DATABASE videm_paren_exists_compare DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_exists_compare
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT ((EXISTS(SELECT 1 FROM t)) = 1) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_exists_compare.t
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
| SHOW CREATE VIEW videm_paren_exists_compare.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select (exists((select 1 from `videm_paren_exists_compare`.`t`)) = 1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_paren_subquery_row_expanded
< ok 0 0
> exec
| CREATE DATABASE videm_paren_subquery_row_expanded DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_paren_subquery_row_expanded
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT (((SELECT 1) = 1) AND (2 = 2)) AS x
< ok 0 0
> query
| SHOW CREATE TABLE videm_paren_subquery_row_expanded.t
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
| SHOW CREATE VIEW videm_paren_subquery_row_expanded.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select (((select 1) = 1) and (2 = 2)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_tsdiff_all_units
< ok 0 0
> exec
| CREATE DATABASE videm_tsdiff_all_units DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_tsdiff_all_units
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME, b DATETIME)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT TIMESTAMPDIFF(MICROSECOND,a,b) AS c1, timestampdiff(second,a,b) AS c2, TIMESTAMPDIFF(MINUTE,a,b) AS c3, TIMESTAMPDIFF(HOUR,a,b) AS c4, TIMESTAMPDIFF(DAY,a,b) AS c5, TIMESTAMPDIFF(WEEK,a,b) AS c6, TIMESTAMPDIFF(MONTH,a,b) AS c7, TIMESTAMPDIFF(QUARTER,a,b) AS c8, TIMESTAMPDIFF(YEAR,a,b) AS c9 FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_tsdiff_all_units.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL,
|   `b` datetime DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_tsdiff_all_units.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select timestampdiff(microsecond,`videm_tsdiff_all_units`.`t`.`a`,`videm_tsdiff_all_units`.`t`.`b`) AS `c1`,timestampdiff(second,`videm_tsdiff_all_units`.`t`.`a`,`videm_tsdiff_all_units`.`t`.`b`) AS `c2`,timestampdiff(minute,`videm_tsdiff_all_units`.`t`.`a`,`videm_tsdiff_all_units`.`t`.`b`) AS `c3`,timestampdiff(hour,`videm_tsdiff_all_units`.`t`.`a`,`videm_tsdiff_all_units`.`t`.`b`) AS `c4`,timestampdiff(day,`videm_tsdiff_all_units`.`t`.`a`,`videm_tsdiff_all_units`.`t`.`b`) AS `c5`,timestampdiff(week,`videm_tsdiff_all_units`.`t`.`a`,`videm_tsdiff_all_units`.`t`.`b`) AS `c6`,timestampdiff(month,`videm_tsdiff_all_units`.`t`.`a`,`videm_tsdiff_all_units`.`t`.`b`) AS `c7`,timestampdiff(quarter,`videm_tsdiff_all_units`.`t`.`a`,`videm_tsdiff_all_units`.`t`.`b`) AS `c8`,timestampdiff(year,`videm_tsdiff_all_units`.`t`.`a`,`videm_tsdiff_all_units`.`t`.`b`) AS `c9` from `videm_tsdiff_all_units`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_tsdiff_column_now
< ok 0 0
> exec
| CREATE DATABASE videm_tsdiff_column_now DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_tsdiff_column_now
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id, TIMESTAMPDIFF(SECOND, created, NOW()) AS age_secs FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_tsdiff_column_now.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_tsdiff_column_now.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `videm_tsdiff_column_now`.`t`.`id` AS `id`,timestampdiff(second,`videm_tsdiff_column_now`.`t`.`created`,now()) AS `age_secs` from `videm_tsdiff_column_now`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_date_arith_family
< ok 0 0
> exec
| CREATE DATABASE videm_date_arith_family DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_date_arith_family
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT TIMESTAMPADD(HOUR, 1, a) AS c1, DATE_ADD(a, INTERVAL 1 DAY) AS c2, DATE_SUB(a, INTERVAL 30 MINUTE) AS c3, ADDDATE(a, INTERVAL 1 WEEK) AS c4, ADDDATE(a, 31) AS c5, SUBDATE(a, INTERVAL 1 MONTH) AS c6, SUBDATE(a, 31) AS c7 FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_date_arith_family.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_date_arith_family.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select timestampadd(HOUR ,1,`videm_date_arith_family`.`t`.`a`) AS `c1`,(`videm_date_arith_family`.`t`.`a` + interval 1 day) AS `c2`,(`videm_date_arith_family`.`t`.`a` - interval 30 minute) AS `c3`,(`videm_date_arith_family`.`t`.`a` + interval 1 week) AS `c4`,(`videm_date_arith_family`.`t`.`a` + interval 31 day) AS `c5`,(`videm_date_arith_family`.`t`.`a` - interval 1 month) AS `c6`,(`videm_date_arith_family`.`t`.`a` - interval 31 day) AS `c7` from `videm_date_arith_family`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_interval_forms
< ok 0 0
> exec
| CREATE DATABASE videm_interval_forms DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_interval_forms
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME, n INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a + INTERVAL 1 DAY AS c1, a - INTERVAL 30 MINUTE AS c2, a + INTERVAL '1:1' DAY_HOUR AS c3, a + INTERVAL '1-1' YEAR_MONTH AS c4, a + INTERVAL (n + 1) DAY AS c5, ADDDATE(a, -5) AS c6, INTERVAL 1 DAY + a AS c7 FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_interval_forms.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL,
|   `n` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_interval_forms.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select (`videm_interval_forms`.`t`.`a` + interval 1 day) AS `c1`,(`videm_interval_forms`.`t`.`a` - interval 30 minute) AS `c2`,(`videm_interval_forms`.`t`.`a` + interval '1:1' day_hour) AS `c3`,(`videm_interval_forms`.`t`.`a` + interval '1-1' year_month) AS `c4`,(`videm_interval_forms`.`t`.`a` + interval (`videm_interval_forms`.`t`.`n` + 1) day) AS `c5`,(`videm_interval_forms`.`t`.`a` + interval -5 day) AS `c6`,(`videm_interval_forms`.`t`.`a` + interval 1 day) AS `c7` from `videm_interval_forms`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_extract_units
< ok 0 0
> exec
| CREATE DATABASE videm_extract_units DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_extract_units
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT EXTRACT(SECOND FROM a) AS c1, extract(day from a) AS c2, EXTRACT(QUARTER FROM a) AS c3, EXTRACT(DAY_HOUR FROM a) AS c4, EXTRACT(YEAR_MONTH FROM a) AS c5, EXTRACT(HOUR_MICROSECOND FROM a) AS c6, EXTRACT(WEEK FROM a) AS c7 FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_extract_units.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_extract_units.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select extract(second from `videm_extract_units`.`t`.`a`) AS `c1`,extract(day from `videm_extract_units`.`t`.`a`) AS `c2`,extract(quarter from `videm_extract_units`.`t`.`a`) AS `c3`,extract(day_hour from `videm_extract_units`.`t`.`a`) AS `c4`,extract(year_month from `videm_extract_units`.`t`.`a`) AS `c5`,extract(hour_microsecond from `videm_extract_units`.`t`.`a`) AS `c6`,extract(week from `videm_extract_units`.`t`.`a`) AS `c7` from `videm_extract_units`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_get_format_kinds
< ok 0 0
> exec
| CREATE DATABASE videm_get_format_kinds DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_get_format_kinds
< ok 0 0
> exec
| CREATE VIEW v AS SELECT GET_FORMAT(DATE, 'USA') AS c1, get_format(time, 'ISO') AS c2, GET_FORMAT(DATETIME, 'EUR') AS c3, GET_FORMAT(TIMESTAMP, 'JIS') AS c4
< ok 0 0
> query
| SHOW CREATE VIEW videm_get_format_kinds.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select get_format(date, 'USA') AS `c1`,get_format(time, 'ISO') AS `c2`,get_format(datetime, 'EUR') AS `c3`,get_format(datetime, 'JIS') AS `c4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_trim_keyword_forms
< ok 0 0
> exec
| CREATE DATABASE videm_trim_keyword_forms DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_trim_keyword_forms
< ok 0 0
> exec
| CREATE TABLE t (s VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT TRIM(BOTH 'x' FROM s) AS c1, TRIM(LEADING 'x' FROM s) AS c2, TRIM(TRAILING 'x' FROM s) AS c3, TRIM('x' FROM s) AS c4, TRIM(s) AS c5 FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_trim_keyword_forms.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `s` varchar(20) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_trim_keyword_forms.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select trim(both 'x' from `videm_trim_keyword_forms`.`t`.`s`) AS `c1`,trim(leading 'x' from `videm_trim_keyword_forms`.`t`.`s`) AS `c2`,trim(trailing 'x' from `videm_trim_keyword_forms`.`t`.`s`) AS `c3`,trim(both 'x' from `videm_trim_keyword_forms`.`t`.`s`) AS `c4`,trim(`videm_trim_keyword_forms`.`t`.`s`) AS `c5` from `videm_trim_keyword_forms`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_substring_keyword_forms
< ok 0 0
> exec
| CREATE DATABASE videm_substring_keyword_forms DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_substring_keyword_forms
< ok 0 0
> exec
| CREATE TABLE t (s VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT SUBSTRING(s FROM 2) AS c1, SUBSTRING(s FROM 2 FOR 3) AS c2, SUBSTRING(s, 2, 3) AS c3 FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_substring_keyword_forms.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `s` varchar(20) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_substring_keyword_forms.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select substr(`videm_substring_keyword_forms`.`t`.`s`,2) AS `c1`,substr(`videm_substring_keyword_forms`.`t`.`s`,2,3) AS `c2`,substr(`videm_substring_keyword_forms`.`t`.`s`,2,3) AS `c3` from `videm_substring_keyword_forms`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS videm_weight_string_forms
< ok 0 0
> exec
| CREATE DATABASE videm_weight_string_forms DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE videm_weight_string_forms
< ok 0 0
> exec
| CREATE TABLE t (s VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT WEIGHT_STRING(s) AS c1, WEIGHT_STRING(s AS CHAR(4)) AS c2, WEIGHT_STRING(s AS BINARY(4)) AS c3 FROM t
< ok 0 0
> query
| SHOW CREATE TABLE videm_weight_string_forms.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `s` varchar(20) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW videm_weight_string_forms.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select weight_string(`videm_weight_string_forms`.`t`.`s`,0,0,0,0) AS `c1`,weight_string(`videm_weight_string_forms`.`t`.`s`,0,4,64,0) AS `c2`,weight_string(`videm_weight_string_forms`.`t`.`s`,0,4,64,1) AS `c3` from `videm_weight_string_forms`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
