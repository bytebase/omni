-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_IndexDiffIdempotence default#1
> exec
| DROP DATABASE IF EXISTS idx_pk_simple
< ok 0 0
> exec
| CREATE DATABASE idx_pk_simple
< ok 1 0
> exec
| USE idx_pk_simple
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE idx_pk_simple.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_pk_composite
< ok 0 0
> exec
| CREATE DATABASE idx_pk_composite
< ok 1 0
> exec
| USE idx_pk_composite
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, c INT, PRIMARY KEY (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE idx_pk_composite.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_unique_single
< ok 0 0
> exec
| CREATE DATABASE idx_unique_single
< ok 1 0
> exec
| USE idx_unique_single
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, UNIQUE KEY ua (a))
< ok 0 0
> query
| SHOW CREATE TABLE idx_unique_single.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_unique_composite
< ok 0 0
> exec
| CREATE DATABASE idx_unique_composite
< ok 1 0
> exec
| USE idx_unique_composite
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, UNIQUE KEY uab (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE idx_unique_composite.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uab` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_secondary
< ok 0 0
> exec
| CREATE DATABASE idx_secondary
< ok 1 0
> exec
| USE idx_secondary
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ka (a), KEY kab (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE idx_secondary.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`),
|   KEY `kab` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_unnamed_autoname
< ok 0 0
> exec
| CREATE DATABASE idx_unnamed_autoname
< ok 1 0
> exec
| USE idx_unnamed_autoname
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY (a), KEY (a,b))
< ok 0 0
> query
| SHOW CREATE TABLE idx_unnamed_autoname.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `a_2` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_using_btree
< ok 0 0
> exec
| CREATE DATABASE idx_using_btree
< ok 1 0
> exec
| USE idx_using_btree
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ua (a) USING BTREE) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE idx_using_btree.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ua` (`a`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_prefix
< ok 0 0
> exec
| CREATE DATABASE idx_prefix
< ok 1 0
> exec
| USE idx_prefix
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100), c CHAR(40), KEY sp (s(10)), KEY cp (c(5)))
< ok 0 0
> query
| SHOW CREATE TABLE idx_prefix.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `s` varchar(100) DEFAULT NULL,
|   `c` char(40) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `sp` (`s`(10)),
|   KEY `cp` (`c`(5))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_comment
< ok 0 0
> exec
| CREATE DATABASE idx_comment
< ok 1 0
> exec
| USE idx_comment
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a) COMMENT 'idx ''c'' here')
< ok 0 0
> query
| SHOW CREATE TABLE idx_comment.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) COMMENT 'idx ''c'' here'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_key_block_size
< ok 0 0
> exec
| CREATE DATABASE idx_key_block_size
< ok 1 0
> exec
| USE idx_key_block_size
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a) KEY_BLOCK_SIZE=4) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE idx_key_block_size.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_descending
< ok 0 0
> exec
| CREATE DATABASE idx_descending
< ok 1 0
> exec
| USE idx_descending
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ab (a, b DESC))
< ok 0 0
> query
| SHOW CREATE TABLE idx_descending.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ab` (`a`,`b` DESC)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_composite_mixed
< ok 0 0
> exec
| CREATE DATABASE idx_composite_mixed
< ok 1 0
> exec
| USE idx_composite_mixed
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, s VARCHAR(50), KEY mix (a DESC, s(8)))
< ok 0 0
> query
| SHOW CREATE TABLE idx_composite_mixed.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `s` varchar(50) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `mix` (`a` DESC,`s`(8))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_invisible
< ok 0 0
> exec
| CREATE DATABASE idx_invisible
< ok 1 0
> exec
| USE idx_invisible
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY va (a) /*!80000 INVISIBLE */)
< ok 0 0
> query
| SHOW CREATE TABLE idx_invisible.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `va` (`a`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_functional
< ok 0 0
> exec
| CREATE DATABASE idx_functional
< ok 1 0
> exec
| USE idx_functional
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY fa ((a + 1)), KEY fab ((a + b)))
< ok 0 0
> query
| SHOW CREATE TABLE idx_functional.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fa` (((`a` + 1))),
|   KEY `fab` (((`a` + `b`)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_mv_unsigned
< ok 0 0
> exec
| CREATE DATABASE idx_mv_unsigned
< ok 1 0
> exec
| USE idx_mv_unsigned
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, j JSON, KEY mv ((CAST(j->'$.ids' AS UNSIGNED ARRAY))))
< ok 0 0
> query
| SHOW CREATE TABLE idx_mv_unsigned.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `mv` ((cast(json_extract(`j`,_utf8mb4'$.ids') as unsigned array)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_mv_char
< ok 0 0
> exec
| CREATE DATABASE idx_mv_char
< ok 1 0
> exec
| USE idx_mv_char
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, j JSON, KEY mv ((CAST(j->'$.codes' AS CHAR(40) ARRAY))))
< ok 0 0
> query
| SHOW CREATE TABLE idx_mv_char.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `mv` ((cast(json_extract(`j`,_utf8mb4'$.codes') as char(40) array)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_mv_decimal
< ok 0 0
> exec
| CREATE DATABASE idx_mv_decimal
< ok 1 0
> exec
| USE idx_mv_decimal
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, j JSON, KEY mv ((CAST(j->'$.amts' AS DECIMAL(10,2) ARRAY))))
< ok 0 0
> query
| SHOW CREATE TABLE idx_mv_decimal.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `mv` ((cast(json_extract(`j`,_utf8mb4'$.amts') as decimal(10, 2) array)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_mv_date
< ok 0 0
> exec
| CREATE DATABASE idx_mv_date
< ok 1 0
> exec
| USE idx_mv_date
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, j JSON, KEY mv ((CAST(j->'$.days' AS DATE ARRAY))))
< ok 0 0
> query
| SHOW CREATE TABLE idx_mv_date.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `mv` ((cast(json_extract(`j`,_utf8mb4'$.days') as date array)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_mv_composite
< ok 0 0
> exec
| CREATE DATABASE idx_mv_composite
< ok 1 0
> exec
| USE idx_mv_composite
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(50), j JSON, KEY mv (name, (CAST(j->'$.ids' AS UNSIGNED ARRAY))))
< ok 0 0
> query
| SHOW CREATE TABLE idx_mv_composite.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `name` varchar(50) DEFAULT NULL,
|   `j` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `mv` (`name`,(cast(json_extract(`j`,_utf8mb4'$.ids') as unsigned array)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_fulltext
< ok 0 0
> exec
| CREATE DATABASE idx_fulltext
< ok 1 0
> exec
| USE idx_fulltext
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, body TEXT, ftcol VARCHAR(200), FULLTEXT KEY ftb (body), FULLTEXT KEY ftc (ftcol))
< ok 0 0
> query
| SHOW CREATE TABLE idx_fulltext.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `body` text,
|   `ftcol` varchar(200) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `ftb` (`body`),
|   FULLTEXT KEY `ftc` (`ftcol`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_spatial
< ok 0 0
> exec
| CREATE DATABASE idx_spatial
< ok 1 0
> exec
| USE idx_spatial
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL, SPATIAL KEY sg (g))
< ok 0 0
> query
| SHOW CREATE TABLE idx_spatial.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `sg` (`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_mixed_all_80
< ok 0 0
> exec
| CREATE DATABASE idx_mixed_all_80
< ok 1 0
> exec
| USE idx_mixed_all_80
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT, b INT, c VARCHAR(40), PRIMARY KEY (id), UNIQUE KEY ua (a), KEY kb (b DESC), KEY kc (c(5)) /*!80000 INVISIBLE */)
< ok 0 0
> query
| SHOW CREATE TABLE idx_mixed_all_80.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `c` varchar(40) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`),
|   KEY `kb` (`b` DESC),
|   KEY `kc` (`c`(5)) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS idx_mixed_all_57
< ok 0 0
> exec
| CREATE DATABASE idx_mixed_all_57
< ok 1 0
> exec
| USE idx_mixed_all_57
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT, b INT, c VARCHAR(40), PRIMARY KEY (id), UNIQUE KEY ua (a), KEY kb (b), KEY kc (c(5)))
< ok 0 0
> query
| SHOW CREATE TABLE idx_mixed_all_57.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `c` varchar(40) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`),
|   KEY `kb` (`b`),
|   KEY `kc` (`c`(5))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
