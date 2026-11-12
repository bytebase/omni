-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTiDBTiFlashLocationLabelsContainer default#1
> exec
| CREATE DATABASE IF NOT EXISTS omni_tf_test
< ok 0 0
> exec
| USE omni_tf_test
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t_tf_l (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE t_tf_l SET TIFLASH REPLICA 0 LOCATION LABELS 'zone'
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_tf_l
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t_tf_ml (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE t_tf_ml SET TIFLASH REPLICA 0 LOCATION LABELS 'zone', 'rack'
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_tf_ml
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t_db_ph (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER DATABASE omni_tf_test SET TIFLASH REPLICA 0
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_db_ph
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t_db_ph2 (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER DATABASE omni_tf_test SET TIFLASH REPLICA 0 LOCATION LABELS 'zone'
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_db_ph2
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_tf_test
< ok 0 0
