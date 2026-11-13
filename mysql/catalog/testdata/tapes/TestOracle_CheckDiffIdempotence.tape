-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_CheckDiffIdempotence default#1
> exec
| DROP DATABASE IF EXISTS chk_unnamed_single
< ok 0 0
> exec
| CREATE DATABASE chk_unnamed_single
< ok 1 0
> exec
| USE chk_unnamed_single
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE chk_unnamed_single.t
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
|   CONSTRAINT `t_chk_1` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_named_single
< ok 0 0
> exec
| CREATE DATABASE chk_named_single
< ok 1 0
> exec
| USE chk_named_single
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT chk_a CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE chk_named_single.t
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
|   CONSTRAINT `chk_a` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_multiple_unnamed
< ok 0 0
> exec
| CREATE DATABASE chk_multiple_unnamed
< ok 1 0
> exec
| USE chk_multiple_unnamed
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a > 0), CHECK (b < 100), CHECK (a < b))
< ok 0 0
> query
| SHOW CREATE TABLE chk_multiple_unnamed.t
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
|   CONSTRAINT `t_chk_1` CHECK ((`a` > 0)),
|   CONSTRAINT `t_chk_2` CHECK ((`b` < 100)),
|   CONSTRAINT `t_chk_3` CHECK ((`a` < `b`))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_mixed_named_unnamed
< ok 0 0
> exec
| CREATE DATABASE chk_mixed_named_unnamed
< ok 1 0
> exec
| USE chk_mixed_named_unnamed
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck1 CHECK (a > 0), CHECK (b < 100))
< ok 0 0
> query
| SHOW CREATE TABLE chk_mixed_named_unnamed.t
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
|   CONSTRAINT `ck1` CHECK ((`a` > 0)),
|   CONSTRAINT `t_chk_1` CHECK ((`b` < 100))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_not_enforced
< ok 0 0
> exec
| CREATE DATABASE chk_not_enforced
< ok 1 0
> exec
| USE chk_not_enforced
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ne CHECK (a < b) NOT ENFORCED)
< ok 0 0
> query
| SHOW CREATE TABLE chk_not_enforced.t
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
|   CONSTRAINT `ne` CHECK ((`a` < `b`)) /*!80016 NOT ENFORCED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_enforced_explicit
< ok 0 0
> exec
| CREATE DATABASE chk_enforced_explicit
< ok 1 0
> exec
| USE chk_enforced_explicit
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT en CHECK (a > 0) ENFORCED)
< ok 0 0
> query
| SHOW CREATE TABLE chk_enforced_explicit.t
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
|   CONSTRAINT `en` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_compound_and
< ok 0 0
> exec
| CREATE DATABASE chk_compound_and
< ok 1 0
> exec
| USE chk_compound_and
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a >= 0 AND b <= 100))
< ok 0 0
> query
| SHOW CREATE TABLE chk_compound_and.t
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
|   CONSTRAINT `t_chk_1` CHECK (((`a` >= 0) and (`b` <= 100)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_compound_or
< ok 0 0
> exec
| CREATE DATABASE chk_compound_or
< ok 1 0
> exec
| USE chk_compound_or
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a = 1 OR b = 2))
< ok 0 0
> query
| SHOW CREATE TABLE chk_compound_or.t
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
|   CONSTRAINT `t_chk_1` CHECK (((`a` = 1) or (`b` = 2)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_between
< ok 0 0
> exec
| CREATE DATABASE chk_between
< ok 1 0
> exec
| USE chk_between
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a BETWEEN 1 AND 10))
< ok 0 0
> query
| SHOW CREATE TABLE chk_between.t
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
|   CONSTRAINT `t_chk_1` CHECK ((`a` between 1 and 10))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_in_list
< ok 0 0
> exec
| CREATE DATABASE chk_in_list
< ok 1 0
> exec
| USE chk_in_list
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a IN (1, 2, 3)))
< ok 0 0
> query
| SHOW CREATE TABLE chk_in_list.t
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
|   CONSTRAINT `t_chk_1` CHECK ((`a` in (1,2,3)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_string_literal
< ok 0 0
> exec
| CREATE DATABASE chk_string_literal
< ok 1 0
> exec
| USE chk_string_literal
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(20), CHECK (s = 'hello')) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE chk_string_literal.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_chk_1` CHECK ((`s` = _utf8mb4'hello'))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_function_call
< ok 0 0
> exec
| CREATE DATABASE chk_function_call
< ok 1 0
> exec
| USE chk_function_call
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(20), CHECK (CHAR_LENGTH(s) > 0))
< ok 0 0
> query
| SHOW CREATE TABLE chk_function_call.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_chk_1` CHECK ((char_length(`s`) > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_backtick_idents
< ok 0 0
> exec
| CREATE DATABASE chk_backtick_idents
< ok 1 0
> exec
| USE chk_backtick_idents
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (`a` < `b`))
< ok 0 0
> query
| SHOW CREATE TABLE chk_backtick_idents.t
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
|   CONSTRAINT `t_chk_1` CHECK ((`a` < `b`))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_odd_spacing
< ok 0 0
> exec
| CREATE DATABASE chk_odd_spacing
< ok 1 0
> exec
| USE chk_odd_spacing
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (   a    >0   ))
< ok 0 0
> query
| SHOW CREATE TABLE chk_odd_spacing.t
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
|   CONSTRAINT `t_chk_1` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_rich_compound
< ok 0 0
> exec
| CREATE DATABASE chk_rich_compound
< ok 1 0
> exec
| USE chk_rich_compound
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, s VARCHAR(20), CONSTRAINT c1 CHECK (UPPER(s) = 'HI' AND a IN (1,2,3)), CONSTRAINT c2 CHECK (a < b) NOT ENFORCED) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE chk_rich_compound.t
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
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `c1` CHECK (((upper(`s`) = _utf8mb4'HI') and (`a` in (1,2,3)))),
|   CONSTRAINT `c2` CHECK ((`a` < `b`)) /*!80016 NOT ENFORCED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS chk_column_level
< ok 0 0
> exec
| CREATE DATABASE chk_column_level
< ok 1 0
> exec
| USE chk_column_level
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE chk_column_level.t
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
|   CONSTRAINT `t_chk_1` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
