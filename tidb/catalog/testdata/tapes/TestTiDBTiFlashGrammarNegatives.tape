-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTiDBTiFlashGrammarNegatives default#1
> exec
| CREATE DATABASE IF NOT EXISTS omni_tfn_test
< ok 0 0
> exec
| USE omni_tfn_test
< ok 0 0
> exec
| CREATE TABLE t_neg (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE t_neg SET TIFLASH REPLICA 0 LOCATION
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 48 near "" 
> exec
| ALTER TABLE t_neg SET TIFLASH REPLICA 0 LOCATION LABELS
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 55 near "" 
> exec
| ALTER TABLE t_neg SET TIFLASH REPLICA 0 LOCATION LABELS 'a',
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 60 near "" 
> exec
| ALTER TABLE t_neg SET TIFLASH REPLICA 0 LOCATION LABELS zone
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 60 near "zone" 
> exec
| CREATE DATABASE omni_tfn_test_x DEFAULT SET TIFLASH REPLICA 1
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 43 near "SET TIFLASH REPLICA 1" 
> exec
| ALTER DATABASE omni_tfn_test DEFAULT SET TIFLASH REPLICA 1
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 40 near "SET TIFLASH REPLICA 1" 
> exec
| DROP DATABASE IF EXISTS omni_tfn_test
< ok 0 0
