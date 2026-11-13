-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestQualifiedEventsAndTriggers default#1
> exec
| DROP DATABASE IF EXISTS a
< ok 0 0
> exec
| DROP DATABASE IF EXISTS b
< ok 0 0
> exec
| CREATE DATABASE a
< ok 1 0
> exec
| CREATE DATABASE b
< ok 1 0
> exec
| USE a
< ok 0 0
> exec
| CREATE TABLE a.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.only_b (id INT)
< ok 0 0
> exec
| CREATE EVENT b.e ON SCHEDULE EVERY 1 DAY DO SELECT 2
< ok 0 0
> exec
| ALTER EVENT b.e COMMENT 'x'
< ok 0 0
> query
| SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')
< result
< column "CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT)" "TEXT"
< row
: bytes
| event b.e x
> exec
| DROP DATABASE IF EXISTS a
< ok 1 0
> exec
| DROP DATABASE IF EXISTS b
< ok 2 0
> exec
| CREATE DATABASE a
< ok 1 0
> exec
| CREATE DATABASE b
< ok 1 0
> exec
| USE a
< ok 0 0
> exec
| CREATE TABLE a.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.only_b (id INT)
< ok 0 0
> exec
| CREATE EVENT a.e ON SCHEDULE EVERY 1 DAY DO SELECT 1
< ok 0 0
> exec
| CREATE EVENT b.e ON SCHEDULE EVERY 1 DAY DO SELECT 2
< ok 0 0
> exec
| DROP EVENT b.e
< ok 0 0
> query
| SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')
< result
< column "CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT)" "TEXT"
< row
: bytes
| event a.e 
> exec
| DROP DATABASE IF EXISTS a
< ok 1 0
> exec
| DROP DATABASE IF EXISTS b
< ok 2 0
> exec
| CREATE DATABASE a
< ok 1 0
> exec
| CREATE DATABASE b
< ok 1 0
> exec
| USE a
< ok 0 0
> exec
| CREATE TABLE a.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.only_b (id INT)
< ok 0 0
> exec
| CREATE EVENT a.e ON SCHEDULE EVERY 1 DAY DO SELECT 1
< ok 0 0
> exec
| CREATE EVENT b.f ON SCHEDULE EVERY 1 DAY DO SELECT 2
< ok 0 0
> exec
| ALTER EVENT b.f RENAME TO a.e
< error
| 1537 HY000 Event 'e' already exists
> query
| SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')
< result
< column "CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT)" "TEXT"
< row
: bytes
| event a.e 
< row
: bytes
| event b.f 
> exec
| DROP DATABASE IF EXISTS a
< ok 1 0
> exec
| DROP DATABASE IF EXISTS b
< ok 2 0
> exec
| CREATE DATABASE a
< ok 1 0
> exec
| CREATE DATABASE b
< ok 1 0
> exec
| USE a
< ok 0 0
> exec
| CREATE TABLE a.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.only_b (id INT)
< ok 0 0
> exec
| CREATE EVENT b.e ON SCHEDULE EVERY 1 DAY DO SELECT 2
< ok 0 0
> exec
| ALTER EVENT b.e RENAME TO f
< ok 0 0
> query
| SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')
< result
< column "CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT)" "TEXT"
< row
: bytes
| event a.f 
> exec
| DROP DATABASE IF EXISTS a
< ok 1 0
> exec
| DROP DATABASE IF EXISTS b
< ok 2 0
> exec
| CREATE DATABASE a
< ok 1 0
> exec
| CREATE DATABASE b
< ok 1 0
> exec
| USE a
< ok 0 0
> exec
| CREATE TABLE a.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.only_b (id INT)
< ok 0 0
> exec
| CREATE EVENT b.e ON SCHEDULE EVERY 1 DAY DO SELECT 2
< ok 0 0
> exec
| ALTER EVENT b.e RENAME TO a.f
< ok 0 0
> query
| SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')
< result
< column "CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT)" "TEXT"
< row
: bytes
| event a.f 
> exec
| DROP DATABASE IF EXISTS a
< ok 1 0
> exec
| DROP DATABASE IF EXISTS b
< ok 2 0
> exec
| CREATE DATABASE a
< ok 1 0
> exec
| CREATE DATABASE b
< ok 1 0
> exec
| USE a
< ok 0 0
> exec
| CREATE TABLE a.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.only_b (id INT)
< ok 0 0
> exec
| CREATE EVENT a.e ON SCHEDULE EVERY 1 DAY DO SELECT 1
< ok 0 0
> exec
| CREATE EVENT a.f ON SCHEDULE EVERY 1 DAY DO SELECT 2
< ok 0 0
> exec
| ALTER EVENT a.e RENAME TO f COMMENT 'new'
< error
| 1537 HY000 Event 'f' already exists
> exec
| ALTER EVENT a.e RENAME TO nosuch.g COMMENT 'new'
< error
| 1049 42000 Unknown database 'nosuch'
> query
| SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')
< result
< column "CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT)" "TEXT"
< row
: bytes
| event a.e 
< row
: bytes
| event a.f 
> exec
| DROP DATABASE IF EXISTS a
< ok 1 0
> exec
| DROP DATABASE IF EXISTS b
< ok 2 0
> exec
| CREATE DATABASE a
< ok 1 0
> exec
| CREATE DATABASE b
< ok 1 0
> exec
| USE a
< ok 0 0
> exec
| CREATE TABLE a.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.only_b (id INT)
< ok 0 0
> exec
| CREATE EVENT a.e ON SCHEDULE EVERY 1 DAY DO SELECT 1
< ok 0 0
> exec
| DROP EVENT b.e
< error
| 1539 HY000 Unknown event 'e'
> exec
| ALTER EVENT b.e COMMENT 'x'
< error
| 1539 HY000 Unknown event 'e'
> query
| SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')
< result
< column "CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT)" "TEXT"
< row
: bytes
| event a.e 
> exec
| DROP DATABASE IF EXISTS a
< ok 1 0
> exec
| DROP DATABASE IF EXISTS b
< ok 2 0
> exec
| CREATE DATABASE a
< ok 1 0
> exec
| CREATE DATABASE b
< ok 1 0
> exec
| USE a
< ok 0 0
> exec
| CREATE TABLE a.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.only_b (id INT)
< ok 0 0
> exec
| CREATE TRIGGER nosuch.tr BEFORE INSERT ON a.t FOR EACH ROW SET NEW.id = 1
< error
| 1435 HY000 Trigger in wrong schema
> query
| SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')
< result
< column "CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT)" "TEXT"
> exec
| DROP DATABASE IF EXISTS a
< ok 1 0
> exec
| DROP DATABASE IF EXISTS b
< ok 2 0
> exec
| CREATE DATABASE a
< ok 1 0
> exec
| CREATE DATABASE b
< ok 1 0
> exec
| USE a
< ok 0 0
> exec
| CREATE TABLE a.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.only_b (id INT)
< ok 0 0
> exec
| CREATE TRIGGER a.tr BEFORE INSERT ON nosuch.t FOR EACH ROW SET NEW.id = 1
< error
| 1049 42000 Unknown database 'nosuch'
> query
| SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')
< result
< column "CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT)" "TEXT"
> exec
| DROP DATABASE IF EXISTS a
< ok 1 0
> exec
| DROP DATABASE IF EXISTS b
< ok 2 0
> exec
| CREATE DATABASE a
< ok 1 0
> exec
| CREATE DATABASE b
< ok 1 0
> exec
| USE a
< ok 0 0
> exec
| CREATE TABLE a.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.only_b (id INT)
< ok 0 0
> exec
| CREATE TRIGGER tr BEFORE INSERT ON only_b FOR EACH ROW SET NEW.id = 1
< error
| 1146 42S02 Table 'a.only_b' doesn't exist
> query
| SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')
< result
< column "CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT)" "TEXT"
> exec
| DROP DATABASE IF EXISTS a
< ok 1 0
> exec
| DROP DATABASE IF EXISTS b
< ok 2 0
> exec
| CREATE DATABASE a
< ok 1 0
> exec
| CREATE DATABASE b
< ok 1 0
> exec
| USE a
< ok 0 0
> exec
| CREATE TABLE a.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.only_b (id INT)
< ok 0 0
> exec
| CREATE TRIGGER b.tr BEFORE INSERT ON a.t FOR EACH ROW SET NEW.id = 1
< error
| 1435 HY000 Trigger in wrong schema
> exec
| CREATE TRIGGER tr BEFORE INSERT ON b.t FOR EACH ROW SET NEW.id = 1
< error
| 1435 HY000 Trigger in wrong schema
> query
| SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')
< result
< column "CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT)" "TEXT"
> exec
| DROP DATABASE IF EXISTS a
< ok 1 0
> exec
| DROP DATABASE IF EXISTS b
< ok 2 0
> exec
| CREATE DATABASE a
< ok 1 0
> exec
| CREATE DATABASE b
< ok 1 0
> exec
| USE a
< ok 0 0
> exec
| CREATE TABLE a.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.only_b (id INT)
< ok 0 0
> exec
| CREATE TRIGGER b.tr BEFORE INSERT ON b.t FOR EACH ROW SET NEW.id = 1
< ok 0 0
> query
| SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')
< result
< column "CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT)" "TEXT"
< row
: bytes
| trigger b.tr on t
> exec
| DROP DATABASE IF EXISTS a
< ok 1 0
> exec
| DROP DATABASE IF EXISTS b
< ok 2 0
> exec
| CREATE DATABASE a
< ok 1 0
> exec
| CREATE DATABASE b
< ok 1 0
> exec
| USE a
< ok 0 0
> exec
| CREATE TABLE a.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.only_b (id INT)
< ok 0 0
> exec
| CREATE TRIGGER nosuch.tr BEFORE INSERT ON t FOR EACH ROW SET NEW.id = 1
< error
| 1049 42000 Unknown database 'nosuch'
> query
| SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')
< result
< column "CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT)" "TEXT"
> exec
| DROP DATABASE IF EXISTS a
< ok 1 0
> exec
| DROP DATABASE IF EXISTS b
< ok 2 0
> exec
| CREATE DATABASE a
< ok 1 0
> exec
| CREATE DATABASE b
< ok 1 0
> exec
| USE a
< ok 0 0
> exec
| CREATE TABLE a.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.only_b (id INT)
< ok 0 0
> exec
| CREATE TRIGGER b.tr BEFORE INSERT ON a.nosuch FOR EACH ROW SET NEW.id = 1
< error
| 1146 42S02 Table 'a.nosuch' doesn't exist
> query
| SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')
< result
< column "CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT)" "TEXT"
> exec
| DROP DATABASE IF EXISTS a
< ok 1 0
> exec
| DROP DATABASE IF EXISTS b
< ok 2 0
> exec
| CREATE DATABASE a
< ok 1 0
> exec
| CREATE DATABASE b
< ok 1 0
> exec
| USE a
< ok 0 0
> exec
| CREATE TABLE a.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.only_b (id INT)
< ok 0 0
> exec
| CREATE EVENT nosuch.e ON SCHEDULE EVERY 1 DAY DO SELECT 1
< error
| 1049 42000 Unknown database 'nosuch'
> query
| SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')
< result
< column "CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT)" "TEXT"
> exec
| DROP DATABASE IF EXISTS a
< ok 1 0
> exec
| DROP DATABASE IF EXISTS b
< ok 2 0
> exec
| CREATE DATABASE a
< ok 1 0
> exec
| CREATE DATABASE b
< ok 1 0
> exec
| USE a
< ok 0 0
> exec
| CREATE TABLE a.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.only_b (id INT)
< ok 0 0
> exec
| CREATE TRIGGER b.tr BEFORE INSERT ON t FOR EACH ROW SET NEW.id = 1
< ok 0 0
> exec
| CREATE TRIGGER b.tr2 BEFORE INSERT ON only_b FOR EACH ROW SET NEW.id = 1
< ok 0 0
> query
| SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')
< result
< column "CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT)" "TEXT"
< row
: bytes
| trigger b.tr on t
< row
: bytes
| trigger b.tr2 on only_b
> exec
| DROP DATABASE IF EXISTS a
< ok 1 0
> exec
| DROP DATABASE IF EXISTS b
< ok 2 0
> exec
| CREATE DATABASE a
< ok 1 0
> exec
| CREATE DATABASE b
< ok 1 0
> exec
| USE a
< ok 0 0
> exec
| CREATE TABLE a.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.t (id INT)
< ok 0 0
> exec
| CREATE TABLE b.only_b (id INT)
< ok 0 0
> exec
| CREATE EVENT a.e ON SCHEDULE EVERY 1 DAY DO SELECT 1
< ok 0 0
> exec
| CREATE EVENT b.e ON SCHEDULE EVERY 1 DAY DO SELECT 2
< ok 0 0
> query
| SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')
< result
< column "CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT)" "TEXT"
< row
: bytes
| event a.e 
< row
: bytes
| event b.e 
