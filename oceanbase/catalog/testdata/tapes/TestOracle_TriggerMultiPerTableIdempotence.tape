-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
| trig_multi
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
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
| TRIGGER au AFTER UPDATE ON t FOR EACH ROW SET @x = NEW.c
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
| TRIGGER bd BEFORE DELETE ON t FOR EACH ROW SET @y = OLD.a
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
| TRIGGER bi1 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
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
| TRIGGER bi2 BEFORE INSERT ON t FOR EACH ROW FOLLOWS bi1 SET NEW.b = 2
