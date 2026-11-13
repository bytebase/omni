-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestTriggerNameFormsOceanBaseOracle default#1
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| CREATE TRIGGER omni_t_trigger_name_forms.n01 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER `omni_t_trigger_name_forms`.`n02` BEFORE INSERT ON `omni_t_trigger_name_forms`.`t` FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER omni_t_trigger_name_forms . n03 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER IF NOT EXISTS omni_t_trigger_name_forms.n04 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER n05 BEFORE INSERT ON omni_t_trigger_name_forms.t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER omni_t_trigger_name_forms.n06 AFTER UPDATE ON t FOR EACH ROW BEGIN SET @x = 1; END
< ok 0 0
> exec
| CREATE TRIGGER omni_t_trigger_name_forms.n07 BEFORE INSERT ON t FOR EACH ROW FOLLOWS omni_t_trigger_name_forms.n01 SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER omni_t_trigger_name_forms.n08 BEFORE INSERT ON t FOR EACH ROW PRECEDES n01 SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER omni_t_trigger_name_forms.n09 BEFORE INSERT ON t FOR EACH ROW FOLLOWS `omni_t_trigger_name_forms`.`n01` SET NEW.a = 1
< ok 0 0
> exec
| CREATE TRIGGER omni_t_no_such_database.n10 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< error
| 1049 42000 Unknown database
> exec
| DROP TRIGGER omni_t_trigger_name_forms.n05
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS omni_t_trigger_name_forms.n11
< ok 0 0
> exec
| CREATE TRIGGER a.omni_t_trigger_name_forms.n12 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< error
| 1064 42000 …near '.' at line 1
> exec
| CREATE TRIGGER omni_t_trigger_name_forms. BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< error
| 1064 42000 …near 'BEFORE' at line 1
> exec
| CREATE TRIGGER .n13 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< error
| 1064 42000 …near '.' at line 1
> exec
| CREATE TRIGGER 'omni_t_trigger_name_forms'.n14 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< error
| 1064 42000 …near ''omni_t_trigger_name_forms'' at line 1
> exec
| CREATE TRIGGER omni_t_trigger_name_forms.'n15' BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< error
| 1064 42000 …near ''n15'' at line 1
> exec
| CREATE TRIGGER omni_t_trigger_name_forms.n16 BEFORE INSERT ON a.omni_t_trigger_name_forms.t FOR EACH ROW SET NEW.a = 1
< error
| 1064 42000 …near '.' at line 1
> exec
| CREATE TRIGGER omni_t_trigger_name_forms.n17 BEFORE INSERT ON t FOR EACH ROW FOLLOWS a.omni_t_trigger_name_forms.n01 SET NEW.a = 1
< error
| 1064 42000 …near '.' at line 1
> exec
| CREATE TRIGGER omni_t_trigger_name_forms.n18 BEFORE INSERT ON t FOR EACH ROW FOLLOWS SET NEW.a = 1
< error
| 1064 42000 …near 'SET' at line 1
> exec
| DROP TRIGGER a.omni_t_trigger_name_forms.n01
< error
| 1064 42000 …near '.' at line 1
