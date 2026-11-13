-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestScenario_C21 default#1
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
| CREATE TABLE t (a INT DEFAULT NULL)
< ok 0 0
> exec
| CREATE TABLE t2 (c INT)
< ok 0 0
> query
| SELECT COLUMN_DEFAULT, IS_NULLABLE FROM information_schema.COLUMNS
|                WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "COLUMN_DEFAULT" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| NULL
: bytes
| YES
> query
| SELECT COLUMN_DEFAULT, IS_NULLABLE FROM information_schema.COLUMNS
|                WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t2' AND COLUMN_NAME='c'
< result
< column "COLUMN_DEFAULT" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| NULL
: bytes
| YES
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
| CREATE TABLE t1 (a INT)
< ok 0 0
> exec
| CREATE TABLE t2 (a INT)
< ok 0 0
> exec
| USE testdb; SELECT * FROM t1 JOIN t2 ON t1.a = t2.a
< ok 0 0
> exec
| USE testdb; SELECT * FROM t1 INNER JOIN t2 ON t1.a = t2.a
< ok 0 0
> exec
| USE testdb; SELECT * FROM t1 CROSS JOIN t2
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
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| USE testdb; SELECT * FROM t ORDER BY a
< ok 0 0
> exec
| USE testdb; SELECT * FROM t ORDER BY a ASC
< ok 0 0
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
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| USE testdb; SELECT * FROM t LIMIT 10
< ok 0 0
> exec
| USE testdb; SELECT * FROM t LIMIT 10 OFFSET 0
< ok 0 0
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
| CREATE TABLE parent (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE child (p INT, FOREIGN KEY (p) REFERENCES parent(id))
< ok 0 0
> query
| SELECT DELETE_RULE, UPDATE_RULE FROM information_schema.REFERENTIAL_CONSTRAINTS
|               WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='child'
< result
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< row
: bytes
| RESTRICT
: bytes
| RESTRICT
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
| CREATE TABLE parent (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE child (p INT, FOREIGN KEY (p) REFERENCES parent(id) ON DELETE CASCADE)
< ok 0 0
> query
| SELECT DELETE_RULE, UPDATE_RULE FROM information_schema.REFERENTIAL_CONSTRAINTS
|               WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='child'
< result
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< row
: bytes
| CASCADE
: bytes
| RESTRICT
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
| CREATE TABLE t (a INT, KEY k (a))
< ok 0 0
> exec
| CREATE TABLE t_b (a INT, KEY kb (a) USING BTREE)
< ok 0 0
> query
| SELECT INDEX_TYPE FROM information_schema.STATISTICS
|               WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND INDEX_NAME='k'
< result
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| BTREE
> query
| SELECT INDEX_TYPE FROM information_schema.STATISTICS
|               WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t_b' AND INDEX_NAME='kb'
< result
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| BTREE
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
| CREATE TABLE t (a INT, b INT, c INT)
< ok 0 0
> exec
| USE testdb; INSERT INTO t VALUES (1, 2, 3)
< ok 1 0
> exec
| USE testdb; INSERT INTO t (a, b, c) VALUES (4, 5, 6)
< ok 1 0
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
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v1 AS SELECT 1 AS x
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED VIEW v2 AS SELECT 1 AS x
< ok 0 0
> query
| SHOW CREATE VIEW v1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v1` AS select 1 AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE VIEW v2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v2` AS select 1 AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 3 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t_noeng (a INT)
< ok 0 0
> exec
| CREATE TABLE t_eng (a INT) ENGINE=InnoDB
< ok 0 0
> query
| SELECT ENGINE FROM information_schema.TABLES
|               WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t_noeng'
< result
< column "ENGINE" "VARCHAR"
< row
: bytes
| InnoDB
> query
| SELECT ENGINE FROM information_schema.TABLES
|               WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t_eng'
< result
< column "ENGINE" "VARCHAR"
< row
: bytes
| InnoDB
