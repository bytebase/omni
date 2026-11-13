-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_TriggerMigrationApplyCorrectness default#1
> exec
| DROP DATABASE IF EXISTS trig_from_create_before_insert
< ok 0 0
> exec
| CREATE DATABASE trig_from_create_before_insert DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_from_create_before_insert
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> query
| SHOW CREATE TABLE trig_from_create_before_insert.t
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
| trig_from_create_before_insert
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
> exec
| DROP DATABASE IF EXISTS trig_to_create_before_insert
< ok 0 0
> exec
| CREATE DATABASE trig_to_create_before_insert DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_to_create_before_insert
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TABLE trig_to_create_before_insert.t
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
| trig_to_create_before_insert
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
| DROP DATABASE IF EXISTS trigdb
< ok 0 0
> exec
| CREATE DATABASE trigdb DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trigdb
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER `trigdb`.`tr` BEFORE INSERT ON `trigdb`.`t` FOR EACH ROW SET NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TABLE trigdb.t
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
| trigdb
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
| TRIGGER `trigdb`.`tr` BEFORE INSERT ON `trigdb`.`t` FOR EACH ROW SET NEW.a = 1
> exec
| DROP DATABASE IF EXISTS trig_from_create_after_update
< ok 0 0
> exec
| CREATE DATABASE trig_from_create_after_update DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_from_create_after_update
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> query
| SHOW CREATE TABLE trig_from_create_after_update.t
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
| trig_from_create_after_update
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
> exec
| DROP DATABASE IF EXISTS trig_to_create_after_update
< ok 0 0
> exec
| CREATE DATABASE trig_to_create_after_update DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_to_create_after_update
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr AFTER UPDATE ON t FOR EACH ROW SET @x = NEW.b
< ok 0 0
> query
| SHOW CREATE TABLE trig_to_create_after_update.t
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
| trig_to_create_after_update
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
| DROP DATABASE IF EXISTS trigdb
< ok 1 0
> exec
| CREATE DATABASE trigdb DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trigdb
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER `trigdb`.`tr` AFTER UPDATE ON `trigdb`.`t` FOR EACH ROW SET @x = NEW.b
< ok 0 0
> query
| SHOW CREATE TABLE trigdb.t
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
| trigdb
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
| TRIGGER `trigdb`.`tr` AFTER UPDATE ON `trigdb`.`t` FOR EACH ROW SET @x = NEW.b
> exec
| DROP DATABASE IF EXISTS trig_from_create_before_delete
< ok 0 0
> exec
| CREATE DATABASE trig_from_create_before_delete DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_from_create_before_delete
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> query
| SHOW CREATE TABLE trig_from_create_before_delete.t
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
| trig_from_create_before_delete
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
> exec
| DROP DATABASE IF EXISTS trig_to_create_before_delete
< ok 0 0
> exec
| CREATE DATABASE trig_to_create_before_delete DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_to_create_before_delete
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE DELETE ON t FOR EACH ROW SET @x = OLD.a
< ok 0 0
> query
| SHOW CREATE TABLE trig_to_create_before_delete.t
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
| trig_to_create_before_delete
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
| DROP DATABASE IF EXISTS trigdb
< ok 1 0
> exec
| CREATE DATABASE trigdb DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trigdb
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER `trigdb`.`tr` BEFORE DELETE ON `trigdb`.`t` FOR EACH ROW SET @x = OLD.a
< ok 0 0
> query
| SHOW CREATE TABLE trigdb.t
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
| trigdb
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
| TRIGGER `trigdb`.`tr` BEFORE DELETE ON `trigdb`.`t` FOR EACH ROW SET @x = OLD.a
> exec
| DROP DATABASE IF EXISTS trig_from_create_begin_end
< ok 0 0
> exec
| CREATE DATABASE trig_from_create_begin_end DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_from_create_begin_end
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> query
| SHOW CREATE TABLE trig_from_create_begin_end.t
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
| trig_from_create_begin_end
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
> exec
| DROP DATABASE IF EXISTS trig_to_create_begin_end
< ok 0 0
> exec
| CREATE DATABASE trig_to_create_begin_end DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_to_create_begin_end
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW BEGIN
|   SET NEW.a = 1;
|   SET NEW.b = 2;
| END
< ok 0 0
> query
| SHOW CREATE TABLE trig_to_create_begin_end.t
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
| trig_to_create_begin_end
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
|   SET NEW.b = 2;
| END
> exec
| DROP DATABASE IF EXISTS trigdb
< ok 1 0
> exec
| CREATE DATABASE trigdb DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trigdb
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER `trigdb`.`tr` BEFORE INSERT ON `trigdb`.`t` FOR EACH ROW BEGIN
|   SET NEW.a = 1;
|   SET NEW.b = 2;
| END
< ok 0 0
> query
| SHOW CREATE TABLE trigdb.t
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
| trigdb
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
| TRIGGER `trigdb`.`tr` BEFORE INSERT ON `trigdb`.`t` FOR EACH ROW BEGIN
|   SET NEW.a = 1;
|   SET NEW.b = 2;
| END
> exec
| DROP DATABASE IF EXISTS trig_from_drop_trigger
< ok 0 0
> exec
| CREATE DATABASE trig_from_drop_trigger DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_from_drop_trigger
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TABLE trig_from_drop_trigger.t
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
| trig_from_drop_trigger
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
| DROP DATABASE IF EXISTS trig_to_drop_trigger
< ok 0 0
> exec
| CREATE DATABASE trig_to_drop_trigger DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_to_drop_trigger
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> query
| SHOW CREATE TABLE trig_to_drop_trigger.t
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
| trig_to_drop_trigger
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
> exec
| DROP DATABASE IF EXISTS trigdb
< ok 1 0
> exec
| CREATE DATABASE trigdb DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trigdb
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| DROP TRIGGER `trigdb`.`tr`
< ok 0 0
> query
| SHOW CREATE TABLE trigdb.t
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
| trigdb
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
> exec
| DROP DATABASE IF EXISTS trig_from_change_body
< ok 0 0
> exec
| CREATE DATABASE trig_from_change_body DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_from_change_body
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TABLE trig_from_change_body.t
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
| trig_from_change_body
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
| DROP DATABASE IF EXISTS trig_to_change_body
< ok 0 0
> exec
| CREATE DATABASE trig_to_change_body DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_to_change_body
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 2
< ok 0 0
> query
| SHOW CREATE TABLE trig_to_change_body.t
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
| trig_to_change_body
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
| TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 2
> exec
| DROP DATABASE IF EXISTS trigdb
< ok 1 0
> exec
| CREATE DATABASE trigdb DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trigdb
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| DROP TRIGGER `trigdb`.`tr`
< ok 0 0
> exec
| CREATE TRIGGER `trigdb`.`tr` BEFORE INSERT ON `trigdb`.`t` FOR EACH ROW SET NEW.a = 2
< ok 0 0
> query
| SHOW CREATE TABLE trigdb.t
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
| trigdb
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
| TRIGGER `trigdb`.`tr` BEFORE INSERT ON `trigdb`.`t` FOR EACH ROW SET NEW.a = 2
> exec
| DROP DATABASE IF EXISTS trig_from_change_timing
< ok 0 0
> exec
| CREATE DATABASE trig_from_change_timing DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_from_change_timing
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET @x = NEW.a
< ok 0 0
> query
| SHOW CREATE TABLE trig_from_change_timing.t
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
| trig_from_change_timing
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
| TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET @x = NEW.a
> exec
| DROP DATABASE IF EXISTS trig_to_change_timing
< ok 0 0
> exec
| CREATE DATABASE trig_to_change_timing DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_to_change_timing
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr AFTER INSERT ON t FOR EACH ROW SET @x = NEW.a
< ok 0 0
> query
| SHOW CREATE TABLE trig_to_change_timing.t
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
| trig_to_change_timing
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
| DROP DATABASE IF EXISTS trigdb
< ok 1 0
> exec
| CREATE DATABASE trigdb DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trigdb
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET @x = NEW.a
< ok 0 0
> exec
| DROP TRIGGER `trigdb`.`tr`
< ok 0 0
> exec
| CREATE TRIGGER `trigdb`.`tr` AFTER INSERT ON `trigdb`.`t` FOR EACH ROW SET @x = NEW.a
< ok 0 0
> query
| SHOW CREATE TABLE trigdb.t
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
| trigdb
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
| TRIGGER `trigdb`.`tr` AFTER INSERT ON `trigdb`.`t` FOR EACH ROW SET @x = NEW.a
> exec
| DROP DATABASE IF EXISTS trig_from_change_event
< ok 0 0
> exec
| CREATE DATABASE trig_from_change_event DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_from_change_event
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TABLE trig_from_change_event.t
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
| trig_from_change_event
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
| DROP DATABASE IF EXISTS trig_to_change_event
< ok 0 0
> exec
| CREATE DATABASE trig_to_change_event DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_to_change_event
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE UPDATE ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TABLE trig_to_change_event.t
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
| trig_to_change_event
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
| TRIGGER tr BEFORE UPDATE ON t FOR EACH ROW SET NEW.a = 1
> exec
| DROP DATABASE IF EXISTS trigdb
< ok 1 0
> exec
| CREATE DATABASE trigdb DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trigdb
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| DROP TRIGGER `trigdb`.`tr`
< ok 0 0
> exec
| CREATE TRIGGER `trigdb`.`tr` BEFORE UPDATE ON `trigdb`.`t` FOR EACH ROW SET NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TABLE trigdb.t
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
| trigdb
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
| TRIGGER `trigdb`.`tr` BEFORE UPDATE ON `trigdb`.`t` FOR EACH ROW SET NEW.a = 1
> exec
| DROP DATABASE IF EXISTS trig_from_multi_add_drop
< ok 0 0
> exec
| CREATE DATABASE trig_from_multi_add_drop DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_from_multi_add_drop
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER keep AFTER UPDATE ON t FOR EACH ROW SET @x = NEW.a
< ok 0 0
> exec
| CREATE TRIGGER goes BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TABLE trig_from_multi_add_drop.t
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
| trig_from_multi_add_drop
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| keep
: bytes
| AFTER
: bytes
| UPDATE
: bytes
| t
: bytes
| TRIGGER keep AFTER UPDATE ON t FOR EACH ROW SET @x = NEW.a
< row
: bytes
| goes
: bytes
| BEFORE
: bytes
| INSERT
: bytes
| t
: bytes
| TRIGGER goes BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
> exec
| DROP DATABASE IF EXISTS trig_to_multi_add_drop
< ok 0 0
> exec
| CREATE DATABASE trig_to_multi_add_drop DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_to_multi_add_drop
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER keep AFTER UPDATE ON t FOR EACH ROW SET @x = NEW.a
< ok 0 0
> exec
| CREATE TRIGGER comes BEFORE DELETE ON t FOR EACH ROW SET @y = OLD.a
< ok 0 0
> query
| SHOW CREATE TABLE trig_to_multi_add_drop.t
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
| trig_to_multi_add_drop
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| keep
: bytes
| AFTER
: bytes
| UPDATE
: bytes
| t
: bytes
| TRIGGER keep AFTER UPDATE ON t FOR EACH ROW SET @x = NEW.a
< row
: bytes
| comes
: bytes
| BEFORE
: bytes
| DELETE
: bytes
| t
: bytes
| TRIGGER comes BEFORE DELETE ON t FOR EACH ROW SET @y = OLD.a
> exec
| DROP DATABASE IF EXISTS trigdb
< ok 1 0
> exec
| CREATE DATABASE trigdb DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trigdb
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TRIGGER keep AFTER UPDATE ON t FOR EACH ROW SET @x = NEW.a
< ok 0 0
> exec
| CREATE TRIGGER goes BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| DROP TRIGGER `trigdb`.`goes`
< ok 0 0
> exec
| CREATE TRIGGER `trigdb`.`comes` BEFORE DELETE ON `trigdb`.`t` FOR EACH ROW SET @y = OLD.a
< ok 0 0
> query
| SHOW CREATE TABLE trigdb.t
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
| trigdb
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| keep
: bytes
| AFTER
: bytes
| UPDATE
: bytes
| t
: bytes
| TRIGGER keep AFTER UPDATE ON t FOR EACH ROW SET @x = NEW.a
< row
: bytes
| comes
: bytes
| BEFORE
: bytes
| DELETE
: bytes
| t
: bytes
| TRIGGER `trigdb`.`comes` BEFORE DELETE ON `trigdb`.`t` FOR EACH ROW SET @y = OLD.a
> exec
| DROP DATABASE IF EXISTS trig_from_modify_table_move_drop_old
< ok 0 0
> exec
| CREATE DATABASE trig_from_modify_table_move_drop_old DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_from_modify_table_move_drop_old
< ok 0 0
> exec
| CREATE TABLE old_t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| CREATE TABLE new_t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON old_t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TABLE trig_from_modify_table_move_drop_old.old_t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| old_t
: bytes
| CREATE TABLE `old_t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE trig_from_modify_table_move_drop_old.new_t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| new_t
: bytes
| CREATE TABLE `new_t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT FROM information_schema.triggers WHERE TRIGGER_SCHEMA = ? ORDER BY EVENT_OBJECT_TABLE, ACTION_TIMING, EVENT_MANIPULATION, ACTION_ORDER
: string
| trig_from_modify_table_move_drop_old
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
| old_t
: bytes
| TRIGGER tr BEFORE INSERT ON old_t FOR EACH ROW SET NEW.a = 1
> exec
| DROP DATABASE IF EXISTS trig_to_modify_table_move_drop_old
< ok 0 0
> exec
| CREATE DATABASE trig_to_modify_table_move_drop_old DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trig_to_modify_table_move_drop_old
< ok 0 0
> exec
| CREATE TABLE new_t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON new_t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TABLE trig_to_modify_table_move_drop_old.new_t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| new_t
: bytes
| CREATE TABLE `new_t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT FROM information_schema.triggers WHERE TRIGGER_SCHEMA = ? ORDER BY EVENT_OBJECT_TABLE, ACTION_TIMING, EVENT_MANIPULATION, ACTION_ORDER
: string
| trig_to_modify_table_move_drop_old
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
| new_t
: bytes
| TRIGGER tr BEFORE INSERT ON new_t FOR EACH ROW SET NEW.a = 1
> exec
| DROP DATABASE IF EXISTS trigdb
< ok 1 0
> exec
| CREATE DATABASE trigdb DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE trigdb
< ok 0 0
> exec
| CREATE TABLE old_t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| CREATE TABLE new_t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON old_t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| DROP TABLE `trigdb`.`old_t`
< ok 0 0
> exec
| CREATE TRIGGER `trigdb`.`tr` BEFORE INSERT ON `trigdb`.`new_t` FOR EACH ROW SET NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TABLE trigdb.new_t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| new_t
: bytes
| CREATE TABLE `new_t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT FROM information_schema.triggers WHERE TRIGGER_SCHEMA = ? ORDER BY EVENT_OBJECT_TABLE, ACTION_TIMING, EVENT_MANIPULATION, ACTION_ORDER
: string
| trigdb
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
| new_t
: bytes
| TRIGGER `trigdb`.`tr` BEFORE INSERT ON `trigdb`.`new_t` FOR EACH ROW SET NEW.a = 1
