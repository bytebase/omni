-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `t_OBCHECK_0000000000000001` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
| 3819 HY000 check constraint violated
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
| 3815 HY000 An expression of a check constraint contains disallowed function.
> exec
| CREATE TABLE t3 (a INT, CHECK (a < NOW()))
< error
| 3814 HY000 An expression of a check constraint contains disallowed function.
> exec
| CREATE TABLE t4 (a INT, CHECK (a < @max))
< error
| 3816 HY000 An expression of a check constraint cannot refer to a user or system variable.
> exec
| CREATE TABLE t5 (a INT, CHECK (a < RAND()))
< error
| 3814 HY000 An expression of a check constraint contains disallowed function.
