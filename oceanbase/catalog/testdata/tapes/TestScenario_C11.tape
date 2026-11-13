-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestScenario_C11 default#1
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
| CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW SET NEW.a = NEW.a
< ok 0 0
> query
| SELECT DEFINER FROM information_schema.TRIGGERS
|              WHERE TRIGGER_SCHEMA='testdb' AND TRIGGER_NAME='trg'
< result
< column "DEFINER" "VARCHAR"
< row
: bytes
| root@%
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
| CREATE DEFINER='root'@'%' TRIGGER trg1 BEFORE INSERT ON t FOR EACH ROW SET NEW.a=1
< ok 0 0
> exec
| CREATE TRIGGER trg2 SQL SECURITY INVOKER BEFORE INSERT ON t FOR EACH ROW SET NEW.a=1
< error
| 1064 42000 …near 'SQL' at line 1
> query
| SELECT COUNT(*) FROM information_schema.COLUMNS
|              WHERE TABLE_SCHEMA='information_schema' AND TABLE_NAME='TRIGGERS'
|              AND COLUMN_NAME='SECURITY_TYPE'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
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
| SET NAMES utf8mb4
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW SET NEW.a = NEW.a
< ok 0 0
> query
| SELECT CHARACTER_SET_CLIENT, COLLATION_CONNECTION, DATABASE_COLLATION
|              FROM information_schema.TRIGGERS
|              WHERE TRIGGER_SCHEMA='testdb' AND TRIGGER_NAME='trg'
< result
< column "CHARACTER_SET_CLIENT" "VARCHAR"
< column "COLLATION_CONNECTION" "VARCHAR"
< column "DATABASE_COLLATION" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
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
| CREATE TRIGGER trg_a BEFORE INSERT ON t FOR EACH ROW SET NEW.a = NEW.a + 1
< ok 0 0
> exec
| CREATE TRIGGER trg_b BEFORE INSERT ON t FOR EACH ROW SET NEW.a = NEW.a + 2
< ok 0 0
> exec
| CREATE TRIGGER trg_c BEFORE INSERT ON t FOR EACH ROW PRECEDES trg_a SET NEW.a = NEW.a + 10
< ok 0 0
> query
| SELECT TRIGGER_NAME, ACTION_ORDER FROM information_schema.TRIGGERS
|              WHERE TRIGGER_SCHEMA='testdb' ORDER BY ACTION_ORDER
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_ORDER" "BIGINT"
< row
: bytes
| trg_c
: int64 1
< row
: bytes
| trg_a
: int64 2
< row
: bytes
| trg_b
: int64 3
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
| CREATE TRIGGER t_ins BEFORE INSERT ON t FOR EACH ROW SET NEW.a = NEW.a + 1
< ok 0 0
> exec
| CREATE TRIGGER t_del BEFORE DELETE ON t FOR EACH ROW SET @x = OLD.a
< ok 0 0
> exec
| CREATE TRIGGER t_upd BEFORE UPDATE ON t FOR EACH ROW SET NEW.a = OLD.a + 1
< ok 0 0
> exec
| CREATE TRIGGER bad1 AFTER INSERT ON t FOR EACH ROW SET @x = OLD.a
< error
| 1363 HY000 There is no OLD row in on INSERT trigger
> exec
| CREATE TRIGGER bad2 AFTER DELETE ON t FOR EACH ROW SET @x = NEW.a
< error
| 1363 HY000 There is no NEW row in on DELETE trigger
> exec
| CREATE TRIGGER bad3 AFTER INSERT ON t FOR EACH ROW SET NEW.a = 99
< error
| 1362 HY000 Updating of NEW row is not allowed in after trigger
> query
| SELECT COUNT(*) FROM information_schema.TRIGGERS
|              WHERE TRIGGER_SCHEMA='testdb' AND TRIGGER_NAME IN ('t_ins','t_del','t_upd')
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 3
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
| CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 4
< ok 0 0
> exec
| CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW SET NEW.b = NEW.a * 2
< ok 0 0
> exec
| ALTER TABLE t COALESCE PARTITION 2
< error
| 1064 42000 …near 'COALESCE PARTITION 2' at line 1
> query
| SELECT COUNT(*) FROM information_schema.TRIGGERS
|              WHERE TRIGGER_SCHEMA='testdb' AND EVENT_OBJECT_TABLE='t' AND TRIGGER_NAME='trg'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 1
