-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
< column "DELETE_RULE" "ENUM"
< column "UPDATE_RULE" "ENUM"
< column "MATCH_OPTION" "ENUM"
< row
: bytes
| NO ACTION
: bytes
| NO ACTION
: bytes
| NONE
> query
| SHOW CREATE TABLE c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `a` int DEFAULT NULL,
|   KEY `a` (`a`),
|   CONSTRAINT `c_ibfk_1` FOREIGN KEY (`a`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
| 1830 HY000 Column 'a' cannot be NOT NULL: needed in a foreign key constraint 'c_ibfk_1' SET NULL
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
< column "MATCH_OPTION" "ENUM"
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
< column "DELETE_RULE" "ENUM"
< column "UPDATE_RULE" "ENUM"
< row
: bytes
| CASCADE
: bytes
| NO ACTION
> query
| SHOW CREATE TABLE c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `a` int DEFAULT NULL,
|   KEY `a` (`a`),
|   CONSTRAINT `c_ibfk_1` FOREIGN KEY (`a`) REFERENCES `p` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
| CREATE TABLE p (id BIGINT PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE c (a INT, FOREIGN KEY (a) REFERENCES p(id))
< error
| 3780 HY000 Referencing column 'a' and referenced column 'id' in foreign key constraint 'c_ibfk_1' are incompatible.
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
< error
| 3733 HY000 Foreign key 'c_ibfk_1' uses virtual column 'b' which is not supported.
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
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   CONSTRAINT `chk_pos` CHECK ((`a` > 0))
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
| CREATE TABLE t1 (a INT CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, CHECK (a > 0))
< ok 0 0
> query
| SELECT CONSTRAINT_NAME, CHECK_CLAUSE
|             FROM information_schema.CHECK_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb' AND CONSTRAINT_NAME LIKE 't1_chk_%'
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t1_chk_1
: bytes
| (`a` > 0)
> query
| SELECT CONSTRAINT_NAME, CHECK_CLAUSE
|             FROM information_schema.CHECK_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb' AND CONSTRAINT_NAME LIKE 't2_chk_%'
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t2_chk_1
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
| 3813 HY000 Column check constraint 't_chk_1' references other column.
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
