-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestScenario_C5 default#1
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
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (a INT, FOREIGN KEY (a) REFERENCES p(id))
< ok 0 0
> query
| SELECT DELETE_RULE, UPDATE_RULE, MATCH_OPTION
|             FROM information_schema.REFERENTIAL_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='c'
< result
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
< row
: bytes
| RESTRICT
: bytes
| RESTRICT
: bytes
| NONE
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE c (
|             a INT NOT NULL,
|             FOREIGN KEY (a) REFERENCES p(id) ON DELETE SET NULL
|         )
< error
| 1830 HY000 Column 'a' cannot be NOT NULL: needed in a foreign key constraint 'c_OBFK_0000000000000001' SET NULL
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
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (a INT, FOREIGN KEY (a) REFERENCES p(id))
< ok 0 0
> query
| SELECT MATCH_OPTION FROM information_schema.REFERENTIAL_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='c'
< result
< column "MATCH_OPTION" "VARCHAR"
< row
: bytes
| NONE
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (a INT, FOREIGN KEY (a) REFERENCES p(id) ON DELETE CASCADE)
< ok 0 0
> query
| SELECT DELETE_RULE, UPDATE_RULE
|             FROM information_schema.REFERENTIAL_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='c'
< result
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< row
: bytes
| CASCADE
: bytes
| RESTRICT
> query
| SHOW CREATE TABLE c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `c_OBFK_0000000000000002` FOREIGN KEY (`a`) REFERENCES `testdb`.`p`(`id`) ON UPDATE RESTRICT ON DELETE CASCADE 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE c (
|             a INT DEFAULT 0,
|             FOREIGN KEY (a) REFERENCES p(id) ON DELETE SET DEFAULT
|         )
< error
| 1215 HY000 Cannot add foreign key constraint
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
| CREATE TABLE p (id BIGINT PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE c (a INT, FOREIGN KEY (a) REFERENCES p(id))
< error
| 1215 HY000 Cannot add foreign key constraint
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
| CREATE TABLE p (id INT PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE c (
|             a INT,
|             b INT GENERATED ALWAYS AS (a+1) VIRTUAL,
|             FOREIGN KEY (b) REFERENCES p(id)
|         )
< ok 0 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
|             a INT,
|             CONSTRAINT chk_pos CHECK (a > 0)
|         )
< ok 0 0
> query
| SELECT ENFORCED FROM information_schema.TABLE_CONSTRAINTS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND CONSTRAINT_NAME='chk_pos'
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
|   CONSTRAINT `chk_pos` CHECK ((`a` > 0))
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
| CREATE TABLE t1 (a INT CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, CHECK (a > 0))
< ok 0 0
> query
| SELECT CONSTRAINT_NAME, CHECK_CLAUSE
|             FROM information_schema.CHECK_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb' AND CONSTRAINT_NAME LIKE 't1_OBCHECK_%'
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t1_OBCHECK_0000000000000003
: bytes
| (`a` > 0)
> query
| SELECT CONSTRAINT_NAME, CHECK_CLAUSE
|             FROM information_schema.CHECK_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb' AND CONSTRAINT_NAME LIKE 't2_OBCHECK_%'
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t2_OBCHECK_0000000000000004
: bytes
| (`a` > 0)
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT CHECK (a > b), b INT)
< error
| 3813 HY000 Column check constraint 't_OBCHECK_0000000000000005' references other column.
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
| CREATE TABLE t (a INT, b INT, CHECK (a > b))
< ok 0 0
