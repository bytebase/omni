-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_QualifiedTriggerDatabase default#1
> exec
| DROP DATABASE IF EXISTS zz_trig_a
< ok 0 0
> exec
| DROP DATABASE IF EXISTS zz_trig_b
< ok 0 0
> exec
| CREATE DATABASE zz_trig_a
< ok 1 0
> exec
| CREATE DATABASE zz_trig_b
< ok 1 0
> exec
| CREATE TABLE zz_trig_a.t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| CREATE TABLE zz_trig_a.only_a (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| CREATE TABLE zz_trig_b.t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| USE zz_trig_a
< ok 0 0
> exec
| CREATE TRIGGER zz_trig_a.q1 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER zz_trig_a.q2 BEFORE INSERT ON zz_trig_a.t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER zz_trig_b.q3 BEFORE INSERT ON zz_trig_b.t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER ZZ_TRIG_A.q4 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER IF NOT EXISTS zz_trig_a.q1 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER zz_trig_b.q5 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< error
| 1435 HY000 Trigger in wrong schema
> exec
| CREATE TRIGGER zz_trig_b.q6 BEFORE INSERT ON zz_trig_a.t FOR EACH ROW SET NEW.a = 1
< error
| 1435 HY000 Trigger in wrong schema
> exec
| CREATE TRIGGER zz_trig_b.q7 BEFORE INSERT ON only_a FOR EACH ROW SET NEW.a = 1
< error
| 1435 HY000 Trigger in wrong schema
> exec
| CREATE TRIGGER zz_trig_a.q8 BEFORE INSERT ON no_such_table FOR EACH ROW SET NEW.a = 1
< error
| 1146 42S02 Table 'zz_trig_a.no_such_table' doesn't exist
> exec
| CREATE TRIGGER zz_trig_b.q9 BEFORE INSERT ON no_such_table FOR EACH ROW SET NEW.a = 1
< error
| 1146 42S02 Table 'zz_trig_a.no_such_table' doesn't exist
> exec
| CREATE TRIGGER zz_no_such_database.q10 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< error
| 1049 42000 Unknown database
> exec
| CREATE TRIGGER zz_no_such_database.q11 BEFORE INSERT ON no_such_table FOR EACH ROW SET NEW.a = 1
< error
| 1049 42000 Unknown database
> exec
| DROP TRIGGER zz_trig_b.q3
< ok 0 0
> query
| SELECT TRIGGER_SCHEMA, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('zz_trig_a', 'zz_trig_b') ORDER BY TRIGGER_SCHEMA, TRIGGER_NAME
< result
< column "TRIGGER_SCHEMA" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| zz_trig_a
: bytes
| q1
< row
: bytes
| zz_trig_a
: bytes
| q2
< row
: bytes
| zz_trig_a
: bytes
| q4
