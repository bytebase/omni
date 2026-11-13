-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_MissedDiffGuards default#1
> exec
| DROP DATABASE IF EXISTS mdg_a
< ok 0 0
> exec
| CREATE DATABASE mdg_a
< ok 1 0
> exec
| USE mdg_a
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10)) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> query
| SHOW CREATE TABLE mdg_a.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` varchar(10) COLLATE utf8mb4_general_ci DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS mdg_b
< ok 0 0
> exec
| CREATE DATABASE mdg_b
< ok 1 0
> exec
| USE mdg_b
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10)) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
< ok 0 0
> query
| SHOW CREATE TABLE mdg_b.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> exec
| DROP DATABASE IF EXISTS mdg_bin
< ok 0 0
> exec
| CREATE DATABASE mdg_bin
< ok 1 0
> exec
| USE mdg_bin
< ok 0 0
> exec
| CREATE TABLE t (b VARCHAR(10) BINARY, p VARCHAR(10)) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE mdg_bin.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   `p` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS mdg_year0
< ok 0 0
> exec
| CREATE DATABASE mdg_year0
< ok 1 0
> exec
| USE mdg_year0
< ok 0 0
> exec
| CREATE TABLE t (n YEAR DEFAULT 0, s YEAR DEFAULT '0')
< ok 0 0
> query
| SHOW CREATE TABLE mdg_year0.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `n` year DEFAULT '0000',
|   `s` year DEFAULT '2000'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
