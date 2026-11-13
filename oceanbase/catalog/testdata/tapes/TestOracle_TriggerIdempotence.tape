-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT FROM information_schema.triggers WHERE TRIGGER_SCHEMA = ? ORDER BY EVENT_OBJECT_TABLE, ACTION_TIMING, EVENT_MANIPULATION, ACTION_ORDER
: string
| trig_idem_before_insert
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
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
| TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
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
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT FROM information_schema.triggers WHERE TRIGGER_SCHEMA = ? ORDER BY EVENT_OBJECT_TABLE, ACTION_TIMING, EVENT_MANIPULATION, ACTION_ORDER
: string
| trig_idem_after_insert
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
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
| TRIGGER tr AFTER INSERT ON t FOR EACH ROW SET @x = NEW.a
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
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT FROM information_schema.triggers WHERE TRIGGER_SCHEMA = ? ORDER BY EVENT_OBJECT_TABLE, ACTION_TIMING, EVENT_MANIPULATION, ACTION_ORDER
: string
| trig_idem_before_update
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
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
| TRIGGER tr BEFORE UPDATE ON t FOR EACH ROW SET NEW.a = OLD.a + 1
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
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT FROM information_schema.triggers WHERE TRIGGER_SCHEMA = ? ORDER BY EVENT_OBJECT_TABLE, ACTION_TIMING, EVENT_MANIPULATION, ACTION_ORDER
: string
| trig_idem_after_update
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
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
| TRIGGER tr AFTER UPDATE ON t FOR EACH ROW SET @x = NEW.b
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
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT FROM information_schema.triggers WHERE TRIGGER_SCHEMA = ? ORDER BY EVENT_OBJECT_TABLE, ACTION_TIMING, EVENT_MANIPULATION, ACTION_ORDER
: string
| trig_idem_before_delete
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
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
| TRIGGER tr BEFORE DELETE ON t FOR EACH ROW SET @x = OLD.a
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
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT FROM information_schema.triggers WHERE TRIGGER_SCHEMA = ? ORDER BY EVENT_OBJECT_TABLE, ACTION_TIMING, EVENT_MANIPULATION, ACTION_ORDER
: string
| trig_idem_after_delete
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
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
| TRIGGER tr AFTER DELETE ON t FOR EACH ROW SET @x = OLD.b
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
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT FROM information_schema.triggers WHERE TRIGGER_SCHEMA = ? ORDER BY EVENT_OBJECT_TABLE, ACTION_TIMING, EVENT_MANIPULATION, ACTION_ORDER
: string
| trig_idem_begin_end
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
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
| TRIGGER tr BEFORE INSERT ON t FOR EACH ROW BEGIN
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
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT FROM information_schema.triggers WHERE TRIGGER_SCHEMA = ? ORDER BY EVENT_OBJECT_TABLE, ACTION_TIMING, EVENT_MANIPULATION, ACTION_ORDER
: string
| trig_idem_string_literal
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
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
| TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.c = (SELECT 'literal value')
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
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT FROM information_schema.triggers WHERE TRIGGER_SCHEMA = ? ORDER BY EVENT_OBJECT_TABLE, ACTION_TIMING, EVENT_MANIPULATION, ACTION_ORDER
: string
| trig_idem_lowercase_body
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
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
| TRIGGER tr before insert on t for each row set new.a = new.b + 1
