-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestScenario_C5 default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
|   KEY `fk_1` (`a`),
|   CONSTRAINT `fk_1` FOREIGN KEY (`a`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
| 1830 HY000 Column 'a' cannot be NOT NULL: needed in a foreign key constraint 'fk_1' SET NULL
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
|   KEY `fk_1` (`a`),
|   CONSTRAINT `fk_1` FOREIGN KEY (`a`) REFERENCES `p` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE p (id BIGINT PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE c (a INT, FOREIGN KEY (a) REFERENCES p(id))
< error
| 3780 HY000 Referencing column 'a' and referenced column 'id' in foreign key constraint 'fk_1' are incompatible.
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
| 3733 HY000 Foreign key 'fk_1' uses virtual column 'b' which is not supported.
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
| SELECT COUNT(*) FROM information_schema.CHECK_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb' AND CONSTRAINT_NAME='chk_pos'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
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
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
> query
| SELECT CONSTRAINT_NAME, CHECK_CLAUSE
|             FROM information_schema.CHECK_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb' AND CONSTRAINT_NAME LIKE 't2_chk_%'
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT CHECK (a > b), b INT)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, CHECK (a > b))
< ok 0 0
