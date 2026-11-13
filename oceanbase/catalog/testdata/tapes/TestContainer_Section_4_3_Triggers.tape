-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_4_3_Triggers default#1
> exec
| DROP TABLE IF EXISTS t_trigger
< ok 0 0
> exec
| CREATE TABLE t_trigger (id INT AUTO_INCREMENT PRIMARY KEY, val INT, name VARCHAR(100))
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS tr_drop_test
< ok 0 0
> exec
| CREATE TRIGGER tr_drop_test BEFORE INSERT ON t_trigger FOR EACH ROW SET NEW.val = 0
< ok 0 0
> exec
| DROP TRIGGER tr_drop_test
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS tr_nonexistent_drop
< ok 0 0
> exec ^
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS tr_no_exist
< ok 0 0
> exec
| DROP TRIGGER tr_no_exist
< error
| 1360 HY000 Trigger does not exist
