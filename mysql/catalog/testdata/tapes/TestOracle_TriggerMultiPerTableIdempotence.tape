-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_TriggerMultiPerTableIdempotence default#1
> exec
| DROP DATABASE IF EXISTS trig_multi
< ok 0 0
> exec
| CREATE DATABASE trig_multi DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_multi
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER bi1 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER bi2 BEFORE INSERT ON t FOR EACH ROW FOLLOWS bi1 SET NEW.b = 2
< ok 0 0
> exec
| CREATE TRIGGER au AFTER UPDATE ON t FOR EACH ROW SET @x = NEW.c
< ok 0 0
> exec
| CREATE TRIGGER bd BEFORE DELETE ON t FOR EACH ROW SET @y = OLD.a
< ok 0 0
> query
| SHOW CREATE TABLE trig_multi.t
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
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT FROM information_schema.triggers WHERE TRIGGER_SCHEMA = ? ORDER BY EVENT_OBJECT_TABLE, ACTION_TIMING, EVENT_MANIPULATION, ACTION_ORDER
: string
| trig_multi
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| bi1
: bytes
| BEFORE
: bytes
| INSERT
: bytes
| t
: bytes
| SET NEW.a = 1
< row
: bytes
| bi2
: bytes
| BEFORE
: bytes
| INSERT
: bytes
| t
: bytes
| SET NEW.b = 2
< row
: bytes
| bd
: bytes
| BEFORE
: bytes
| DELETE
: bytes
| t
: bytes
| SET @y = OLD.a
< row
: bytes
| au
: bytes
| AFTER
: bytes
| UPDATE
: bytes
| t
: bytes
| SET @x = NEW.c
