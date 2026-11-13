-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_TriggerIdempotence default#1
> exec
| DROP DATABASE IF EXISTS trig_idem_before_insert
< ok 0 0
> exec
| CREATE DATABASE trig_idem_before_insert DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_idem_before_insert
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TABLE trig_idem_before_insert.t
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
| trig_idem_before_insert
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| tr
: bytes
| BEFORE
: bytes
| INSERT
: bytes
| t
: bytes
| SET NEW.a = 1
> exec
| DROP DATABASE IF EXISTS trig_idem_after_insert
< ok 0 0
> exec
| CREATE DATABASE trig_idem_after_insert DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_idem_after_insert
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr AFTER INSERT ON t FOR EACH ROW SET @x = NEW.a
< ok 0 0
> query
| SHOW CREATE TABLE trig_idem_after_insert.t
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
| trig_idem_after_insert
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| tr
: bytes
| AFTER
: bytes
| INSERT
: bytes
| t
: bytes
| SET @x = NEW.a
> exec
| DROP DATABASE IF EXISTS trig_idem_before_update
< ok 0 0
> exec
| CREATE DATABASE trig_idem_before_update DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_idem_before_update
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE UPDATE ON t FOR EACH ROW SET NEW.a = OLD.a + 1
< ok 0 0
> query
| SHOW CREATE TABLE trig_idem_before_update.t
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
| trig_idem_before_update
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| tr
: bytes
| BEFORE
: bytes
| UPDATE
: bytes
| t
: bytes
| SET NEW.a = OLD.a + 1
> exec
| DROP DATABASE IF EXISTS trig_idem_after_update
< ok 0 0
> exec
| CREATE DATABASE trig_idem_after_update DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_idem_after_update
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr AFTER UPDATE ON t FOR EACH ROW SET @x = NEW.b
< ok 0 0
> query
| SHOW CREATE TABLE trig_idem_after_update.t
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
| trig_idem_after_update
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| tr
: bytes
| AFTER
: bytes
| UPDATE
: bytes
| t
: bytes
| SET @x = NEW.b
> exec
| DROP DATABASE IF EXISTS trig_idem_before_delete
< ok 0 0
> exec
| CREATE DATABASE trig_idem_before_delete DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_idem_before_delete
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE DELETE ON t FOR EACH ROW SET @x = OLD.a
< ok 0 0
> query
| SHOW CREATE TABLE trig_idem_before_delete.t
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
| trig_idem_before_delete
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| tr
: bytes
| BEFORE
: bytes
| DELETE
: bytes
| t
: bytes
| SET @x = OLD.a
> exec
| DROP DATABASE IF EXISTS trig_idem_after_delete
< ok 0 0
> exec
| CREATE DATABASE trig_idem_after_delete DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_idem_after_delete
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr AFTER DELETE ON t FOR EACH ROW SET @x = OLD.b
< ok 0 0
> query
| SHOW CREATE TABLE trig_idem_after_delete.t
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
| trig_idem_after_delete
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| tr
: bytes
| AFTER
: bytes
| DELETE
: bytes
| t
: bytes
| SET @x = OLD.b
> exec
| DROP DATABASE IF EXISTS trig_idem_begin_end
< ok 0 0
> exec
| CREATE DATABASE trig_idem_begin_end DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_idem_begin_end
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW BEGIN
|   SET NEW.a = 1;
|   IF NEW.b IS NULL THEN
|     SET NEW.b = 0;
|   END IF;
| END
< ok 0 0
> query
| SHOW CREATE TABLE trig_idem_begin_end.t
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
| trig_idem_begin_end
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| tr
: bytes
| BEFORE
: bytes
| INSERT
: bytes
| t
: bytes
| BEGIN
|   SET NEW.a = 1;
|   IF NEW.b IS NULL THEN
|     SET NEW.b = 0;
|   END IF;
| END
> exec
| DROP DATABASE IF EXISTS trig_idem_string_literal
< ok 0 0
> exec
| CREATE DATABASE trig_idem_string_literal DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_idem_string_literal
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.c = (SELECT 'literal value')
< ok 0 0
> query
| SHOW CREATE TABLE trig_idem_string_literal.t
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
| trig_idem_string_literal
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| tr
: bytes
| BEFORE
: bytes
| INSERT
: bytes
| t
: bytes
| SET NEW.c = (SELECT 'literal value')
> exec
| DROP DATABASE IF EXISTS trig_idem_lowercase_body
< ok 0 0
> exec
| CREATE DATABASE trig_idem_lowercase_body DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_idem_lowercase_body
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr before insert on t for each row set new.a = new.b + 1
< ok 0 0
> query
| SHOW CREATE TABLE trig_idem_lowercase_body.t
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
| trig_idem_lowercase_body
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "ENUM"
< column "EVENT_MANIPULATION" "ENUM"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| tr
: bytes
| BEFORE
: bytes
| INSERT
: bytes
| t
: bytes
| set new.a = new.b + 1
