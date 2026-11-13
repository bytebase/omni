-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestScenario_C14 default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, CHECK (a > 0))
< ok 0 0
> query
| SELECT ENFORCED FROM information_schema.TABLE_CONSTRAINTS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND CONSTRAINT_TYPE='CHECK'
< result
< column "ENFORCED" "VARCHAR"
< row
: bytes
| YES
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   CONSTRAINT `t_chk_1` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT c_pos CHECK (a > 0))
< ok 0 0
> exec
| ALTER TABLE t ALTER CHECK c_pos NOT ENFORCED
< ok 0 0
> query
| SELECT ENFORCED FROM information_schema.TABLE_CONSTRAINTS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND CONSTRAINT_NAME='c_pos'
< result
< column "ENFORCED" "VARCHAR"
< row
: bytes
| NO
> exec
| ALTER TABLE t ALTER CHECK c_pos ENFORCED
< ok 0 0
> query
| SELECT ENFORCED FROM information_schema.TABLE_CONSTRAINTS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND CONSTRAINT_NAME='c_pos'
< result
< column "ENFORCED" "VARCHAR"
< row
: bytes
| YES
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
|     a INT,
|     g INT AS (a * 2) STORED,
|     CONSTRAINT c_g_nonneg CHECK (g >= 0)
| )
< ok 0 0
> query
| SELECT CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb' AND CONSTRAINT_NAME='c_g_nonneg'
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| c_g_nonneg
: bytes
| (`g` >= 0)
> exec
| INSERT INTO t (a) VALUES (5)
< ok 1 0
> exec
| INSERT INTO t (a) VALUES (-1)
< error
| 3819 HY000 Check constraint 'c_g_nonneg' is violated.
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t1 (a INT, CHECK (CHAR_LENGTH(CAST(a AS CHAR)) < 10))
< ok 0 0
> exec
| CREATE TABLE other (id INT)
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, CHECK (a IN (SELECT id FROM other)))
< error
| 3815 HY000 An expression of a check constraint 't2_chk_1' contains disallowed function.
> exec
| CREATE TABLE t3 (a INT, CHECK (a < NOW()))
< error
| 3814 HY000 An expression of a check constraint 't3_chk_1' contains disallowed function: now.
> exec
| CREATE TABLE t4 (a INT, CHECK (a < @max))
< error
| 3816 HY000 An expression of a check constraint 't4_chk_1' cannot refer to a user or system variable.
> exec
| CREATE TABLE t5 (a INT, CHECK (a < RAND()))
< error
| 3814 HY000 An expression of a check constraint 't5_chk_1' contains disallowed function: rand.
