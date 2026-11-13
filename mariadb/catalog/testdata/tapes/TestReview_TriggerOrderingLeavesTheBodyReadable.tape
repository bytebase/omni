-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_TriggerOrderingLeavesTheBodyReadable default#1
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
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON t FOR EACH ROW -- keep order
| FOLLOWS tr1 SET NEW.a = 2
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
| CREATE TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON t FOR EACH ROW # keep order
| FOLLOWS tr1 SET NEW.a = 2
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
| CREATE TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON t FOR EACH ROW -- keep order
| PRECEDES tr1 SET NEW.a = 2
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
| CREATE TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON t FOR EACH ROW
| FOLLOWS tr1
| SET NEW.a = 2
< ok 0 0
> query
| SHOW CREATE TRIGGER tr2
< result
< column "Trigger" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< column "Created" "DATETIME"
< row
: bytes
| tr2
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER tr2 BEFORE INSERT ON t FOR EACH ROW
| SET NEW.a = 2
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: time 2026-11-13T00:00:00Z
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
| CREATE TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER tr2 BEFORE INSERT ON t FOR EACH ROW FOLLOWS tr1 SET NEW.a = 2
< ok 0 0
> query
| SHOW CREATE TRIGGER tr2
< result
< column "Trigger" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< column "Created" "DATETIME"
< row
: bytes
| tr2
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER tr2 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 2
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: time 2026-11-13T00:00:00Z
